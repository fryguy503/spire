/* tslint:disable */
/* eslint-disable */
/**
 * Spire
 * Spire API documentation
 *
 * The version of the OpenAPI document: 3.0
 * Contact: akkadius1@gmail.com
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


import { ModelsGridEntry } from './models-grid-entry';
import { ModelsZone } from './models-zone';

/**
 * 
 * @export
 * @interface ModelsGrid
 */
export interface ModelsGrid {
    /**
     * 
     * @type {Array<ModelsGridEntry>}
     * @memberof ModelsGrid
     */
    grid_entries?: Array<ModelsGridEntry>;
    /**
     * 
     * @type {number}
     * @memberof ModelsGrid
     */
    id?: number;
    /**
     * 
     * @type {number}
     * @memberof ModelsGrid
     */
    type?: number;
    /**
     * 
     * @type {number}
     * @memberof ModelsGrid
     */
    type_2?: number;
    /**
     * 
     * @type {ModelsZone}
     * @memberof ModelsGrid
     */
    zone?: ModelsZone;
    /**
     * 
     * @type {number}
     * @memberof ModelsGrid
     */
    zoneid?: number;
}


