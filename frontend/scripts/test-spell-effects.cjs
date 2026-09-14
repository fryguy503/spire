const assert = require('node:assert/strict');
const fs = require('node:fs');
const path = require('node:path');
const vm = require('node:vm');
const { test } = require('node:test');
const ts = require('typescript');

const src = path.resolve(__dirname, '../src');
const spellSource = fs.readFileSync(path.join(src, 'app/spells.ts'), 'utf8');
const modules = new Map();

// Execute the real renderer and constants without initializing browser auth or APIs.
// Unexpected API access fails the test; linked spell data is seeded below.
const unavailable = new Proxy({}, {
  get(_, name) { throw new Error(`Unexpected dependency access: ${String(name)}`); },
});
function loadSource(file) {
  if (modules.has(file)) return modules.get(file);
  const module = { exports: {} };
  modules.set(file, module.exports);
  const code = ts.transpileModule(fs.readFileSync(file, 'utf8'), {
    compilerOptions: {
      module: ts.ModuleKind.CommonJS,
      target: ts.ScriptTarget.ES2020,
      esModuleInterop: true,
    },
  }).outputText;
  const load = name => {
    if (name === 'util') return require('node:util');
    if (name.startsWith('@/app/constants/')) return loadSource(path.join(src, name.slice(2) + '.ts'));
    if (name === '@/constants/app') return { App: { DEBUG: false } };
    return unavailable;
  };
  vm.runInThisContext(`(function(require, module, exports) {${code}\n})`, { filename: file })(load, module, module.exports);
  return module.exports;
}

const { Spells } = loadSource(path.join(src, 'app/spells.ts'));
const { DB_SPA } = loadSource(path.join(src, 'app/constants/eq-spell-constants.ts'));
const { SPELL_SPA_DEFINITIONS } = loadSource(path.join(src, 'app/constants/eq-spell-spa-definitions.ts'));

function makeSpell() {
  const spell = { id: 4505, name: 'Deadeye Discipline', buffduration: 10, targettype: 6 };
  for (let i = 1; i <= 16; i++) spell[`classes_${i}`] = i === 9 ? 54 : 255;
  for (let slot = 1; slot <= 12; slot++) setEffect(spell, slot, 254, 0);
  return spell;
}

function setEffect(spell, slot, id, base, limit = 0, max = 0, formula = 100) {
  Object.assign(spell, {
    [`effectid_${slot}`]: id,
    [`effect_base_value_${slot}`]: base,
    [`effect_limit_value_${slot}`]: limit,
    [`max_${slot}`]: max,
    [`formula_${slot}`]: formula,
  });
}

function rendererCases() {
  const source = ts.createSourceFile('spells.ts', spellSource, ts.ScriptTarget.Latest, true);
  const spells = source.statements.find(node => ts.isClassDeclaration(node) && node.name.text === 'Spells');
  const method = spells.members.find(node => node.name?.getText() === 'getSpellEffectInfo');
  const switches = [];
  function visit(node) {
    if (ts.isSwitchStatement(node) && node.expression.getText().includes('effectid_')) switches.push(node);
    ts.forEachChild(node, visit);
  }
  visit(method);
  assert.equal(switches.length, 1, 'Locate the effect ID switch, excluding formula and other nested switches');
  return switches[0].caseBlock.clauses.filter(ts.isCaseClause);
}

test('every catalog effect except the unused-slot sentinel has a renderer case', t => {
  const ids = new Set(rendererCases().map(node => Number(node.expression.getText())));
  const catalog = new Set([...Object.keys(DB_SPA), ...Object.keys(SPELL_SPA_DEFINITIONS)].map(Number));
  const missing = [...catalog].filter(id => id !== 254 && !ids.has(id));
  assert.deepEqual(missing, [], `Missing renderers: ${missing.map(id => `${id}: ${DB_SPA[id]}`).join(', ')}`);
  t.diagnostic(`Scanned ${catalog.size} catalog effects; 254 is the unused-slot sentinel.`);
});

test('effect IDs have no duplicate or accidentally empty cases', () => {
  const seen = new Set();
  for (const node of rendererCases()) {
    const id = Number(node.expression.getText());
    assert.ok(!seen.has(id), `Duplicate renderer for ${id}: ${DB_SPA[id]}`);
    seen.add(id);
    // 146 stores portal coordinates and is explicitly unimplemented in the catalog.
    if (id !== 146) {
      assert.ok(node.statements.some(statement => !ts.isBreakStatement(statement)), `Empty renderer for ${id}`);
    }
  }
});

test('Deadeye Discipline keeps all seven effects and distinguishes critical chance from damage', async () => {
  const spell = makeSpell();
  const slots = [[200, 200, -1], [170, 1000, -1], [155, 100, -1], [273, 1000, -1], [375, 100, -1], [85, 1005, 100], [523, -1500, 0, 1250]];
  slots.forEach((effect, index) => setEffect(spell, index + 1, ...effect));
  Spells.setSpell(1005, { id: 1005, name: 'Deadly Poison', targettype: 5, new_icon: 0 });
  const actual = await Promise.all(slots.map((_, index) => Spells.getSpellEffectInfo(spell, index + 1)));
  const text = actual.map(result => result.info.replace(/<[^>]*>/g, '').replace(/\s+/g, ' ').trim());
  assert.deepEqual(text, [
    '1) Increase Worn Proc Rate by 200%',
    '2) Increase Chance to Critical Nuke by 1000%',
    '3) Increase Critical Nuke Damage by 100% of Base Damage',
    '4) Increase Chance to Critical DoT by 1000%',
    '5) Increase Critical DoT Damage by 100% of Base Damage',
    '6) Add Melee Proc Deadly Poison with 100 % Rate Mod',
    '7) Decrease Current Endurance by 15% up to 1250',
  ]);
  assert.deepEqual(actual.map(result => result.index), [1, 2, 3, 4, 5, 6, 7]);
});

for (const [id, label, suffix] of [
  [155, 'Critical Nuke Damage', '% of Base Damage'],
  [170, 'Chance to Critical Nuke', '%'],
  [508, 'Spell Power', ' (Focus Spell DOT, DD and Healing)'],
]) {
  test(`effect ${id} renders signed and zero values in different slots, ignoring unused limits`, async () => {
    for (const slot of [1, 3, 12]) {
      for (const base of [-100, 0, 100]) {
        for (const limit of [-1, 0, 1]) {
          const spell = makeSpell();
          setEffect(spell, slot, id, base, limit);
          const result = await Spells.getSpellEffectInfo(spell, slot);
          assert.deepEqual(result, {
            index: slot,
            info: `${slot}) ${base < 0 ? 'Decrease' : 'Increase'} ${label} by ${Math.abs(base)}${suffix}`,
          });
        }
      }
    }
  });
}

test('effect 504 renders front and rear flat damage without dividing the amount by ten', async () => {
  for (const [limit, arc] of [[0, 'Rear'], [1, 'Frontal']]) {
    for (const base of [-123, 0, 123]) {
      const spell = makeSpell();
      setEffect(spell, 12, 504, base, limit);
      const result = await Spells.getSpellEffectInfo(spell, 12);
      assert.equal(result.info, `12) ${base < 0 ? 'Decrease' : 'Increase'} ${arc} Arc Melee Damage Amount by ${Math.abs(base)}`);
    }
  }
});

test('new damage renderers retain formula scaling and caps', async () => {
  for (const [id, label, suffix] of [
    [155, 'Critical Nuke Damage', '% of Base Damage'],
    [504, 'Rear Arc Melee Damage Amount', ''],
  ]) {
    const spell = makeSpell();
    setEffect(spell, 3, id, 100, 0, 180, 102);
    const scaled = await Spells.getSpellEffectInfo(spell, 3);
    const unit = id === 155 ? '%' : '';
    const tail = id === 155 ? ' of Base Damage' : suffix;
    assert.equal(scaled.info, `3) Increase ${label} by 154${unit} (L54) to 180${unit} (L80)${tail}`);
    setEffect(spell, 3, id, 100, 0, 80);
    const capped = await Spells.getSpellEffectInfo(spell, 3);
    assert.equal(capped.info, `3) Increase ${label} by 80${suffix}`);
  }
});

test('effect 508 uses the raw focus amount regardless of Formula and Max', async () => {
  for (const base of [-100, 0, 100]) {
    for (const [formula, max] of [[100, 80], [102, 180], [102, 0]]) {
      const spell = makeSpell();
      setEffect(spell, 3, 508, base, 0, max, formula);
      const result = await Spells.getSpellEffectInfo(spell, 3);
      assert.equal(result.info, `3) ${base < 0 ? 'Decrease' : 'Increase'} Spell Power by ${Math.abs(base)} (Focus Spell DOT, DD and Healing)`);
    }
  }
});

test('unused slots and the zero-CHA placeholder stay invisible', async () => {
  const spell = makeSpell();
  assert.equal((await Spells.getSpellEffectInfo(spell, 12)).info, '');
  setEffect(spell, 3, 10, 0);
  assert.equal((await Spells.getSpellEffectInfo(spell, 3)).info, '');
});
