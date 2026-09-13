import type { ModelsUser } from '@/app/api/models/models-user'
import { AppEnv } from '@/app/env/app-env'

// This controls navigation only. The API enforces connection permissions for
// each resource; the instance-admin flag is for managing Spire users.
export function canAccessServerAdmin (user: Pick<ModelsUser, 'id' | 'is_admin'> | null): boolean {
  if (typeof AppEnv.getEnv() === 'undefined') {
    return false
  }

  if (!AppEnv.isLocalAuthEnabled() && !AppEnv.isGithubAuthEnabled()) {
    return true
  }

  return Boolean(user && user.id && (AppEnv.isAppLocal() || user.is_admin))
}
