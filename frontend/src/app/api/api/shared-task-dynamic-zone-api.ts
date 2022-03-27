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


import globalAxios, { AxiosPromise, AxiosInstance } from 'axios';
import { Configuration } from '../configuration';
// Some imports not used depending on template conditions
// @ts-ignore
import { BASE_PATH, COLLECTION_FORMATS, RequestArgs, BaseAPI, RequiredError } from '../base';
// @ts-ignore
import { CrudcontrollersBulkFetchByIdsGetRequest } from '../models';
// @ts-ignore
import { ModelsSharedTaskDynamicZone } from '../models';
/**
 * SharedTaskDynamicZoneApi - axios parameter creator
 * @export
 */
export const SharedTaskDynamicZoneApiAxiosParamCreator = function (configuration?: Configuration) {
    return {
        /**
         * 
         * @summary Creates SharedTaskDynamicZone
         * @param {ModelsSharedTaskDynamicZone} sharedTaskDynamicZone SharedTaskDynamicZone
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        createSharedTaskDynamicZone: async (sharedTaskDynamicZone: ModelsSharedTaskDynamicZone, options: any = {}): Promise<RequestArgs> => {
            // verify required parameter 'sharedTaskDynamicZone' is not null or undefined
            if (sharedTaskDynamicZone === null || sharedTaskDynamicZone === undefined) {
                throw new RequiredError('sharedTaskDynamicZone','Required parameter sharedTaskDynamicZone was null or undefined when calling createSharedTaskDynamicZone.');
            }
            const localVarPath = `/shared_task_dynamic_zone`;
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, 'https://example.com');
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'PUT', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;


    
            localVarHeaderParameter['Content-Type'] = 'application/json';

            const queryParameters = new URLSearchParams(localVarUrlObj.search);
            for (const key in localVarQueryParameter) {
                queryParameters.set(key, localVarQueryParameter[key]);
            }
            for (const key in options.query) {
                queryParameters.set(key, options.query[key]);
            }
            localVarUrlObj.search = (new URLSearchParams(queryParameters)).toString();
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};
            const nonString = typeof sharedTaskDynamicZone !== 'string';
            const needsSerialization = nonString && configuration && configuration.isJsonMime
                ? configuration.isJsonMime(localVarRequestOptions.headers['Content-Type'])
                : nonString;
            localVarRequestOptions.data =  needsSerialization
                ? JSON.stringify(sharedTaskDynamicZone !== undefined ? sharedTaskDynamicZone : {})
                : (sharedTaskDynamicZone || "");

            return {
                url: localVarUrlObj.pathname + localVarUrlObj.search + localVarUrlObj.hash,
                options: localVarRequestOptions,
            };
        },
        /**
         * 
         * @summary Deletes SharedTaskDynamicZone
         * @param {number} id sharedTaskId
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        deleteSharedTaskDynamicZone: async (id: number, options: any = {}): Promise<RequestArgs> => {
            // verify required parameter 'id' is not null or undefined
            if (id === null || id === undefined) {
                throw new RequiredError('id','Required parameter id was null or undefined when calling deleteSharedTaskDynamicZone.');
            }
            const localVarPath = `/shared_task_dynamic_zone/{id}`
                .replace(`{${"id"}}`, encodeURIComponent(String(id)));
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, 'https://example.com');
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'DELETE', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;


    
            const queryParameters = new URLSearchParams(localVarUrlObj.search);
            for (const key in localVarQueryParameter) {
                queryParameters.set(key, localVarQueryParameter[key]);
            }
            for (const key in options.query) {
                queryParameters.set(key, options.query[key]);
            }
            localVarUrlObj.search = (new URLSearchParams(queryParameters)).toString();
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};

            return {
                url: localVarUrlObj.pathname + localVarUrlObj.search + localVarUrlObj.hash,
                options: localVarRequestOptions,
            };
        },
        /**
         * 
         * @summary Gets SharedTaskDynamicZone
         * @param {number} id Id
         * @param {string} [includes] Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names 
         * @param {string} [select] Column names [.] separated to fetch specific fields in response
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getSharedTaskDynamicZone: async (id: number, includes?: string, select?: string, options: any = {}): Promise<RequestArgs> => {
            // verify required parameter 'id' is not null or undefined
            if (id === null || id === undefined) {
                throw new RequiredError('id','Required parameter id was null or undefined when calling getSharedTaskDynamicZone.');
            }
            const localVarPath = `/shared_task_dynamic_zone/{id}`
                .replace(`{${"id"}}`, encodeURIComponent(String(id)));
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, 'https://example.com');
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'GET', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;

            if (includes !== undefined) {
                localVarQueryParameter['includes'] = includes;
            }

            if (select !== undefined) {
                localVarQueryParameter['select'] = select;
            }


    
            const queryParameters = new URLSearchParams(localVarUrlObj.search);
            for (const key in localVarQueryParameter) {
                queryParameters.set(key, localVarQueryParameter[key]);
            }
            for (const key in options.query) {
                queryParameters.set(key, options.query[key]);
            }
            localVarUrlObj.search = (new URLSearchParams(queryParameters)).toString();
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};

            return {
                url: localVarUrlObj.pathname + localVarUrlObj.search + localVarUrlObj.hash,
                options: localVarRequestOptions,
            };
        },
        /**
         * 
         * @summary Gets SharedTaskDynamicZones in bulk
         * @param {CrudcontrollersBulkFetchByIdsGetRequest} body body
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getSharedTaskDynamicZonesBulk: async (body: CrudcontrollersBulkFetchByIdsGetRequest, options: any = {}): Promise<RequestArgs> => {
            // verify required parameter 'body' is not null or undefined
            if (body === null || body === undefined) {
                throw new RequiredError('body','Required parameter body was null or undefined when calling getSharedTaskDynamicZonesBulk.');
            }
            const localVarPath = `/shared_task_dynamic_zones/bulk`;
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, 'https://example.com');
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'POST', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;


    
            localVarHeaderParameter['Content-Type'] = 'application/json';

            const queryParameters = new URLSearchParams(localVarUrlObj.search);
            for (const key in localVarQueryParameter) {
                queryParameters.set(key, localVarQueryParameter[key]);
            }
            for (const key in options.query) {
                queryParameters.set(key, options.query[key]);
            }
            localVarUrlObj.search = (new URLSearchParams(queryParameters)).toString();
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};
            const nonString = typeof body !== 'string';
            const needsSerialization = nonString && configuration && configuration.isJsonMime
                ? configuration.isJsonMime(localVarRequestOptions.headers['Content-Type'])
                : nonString;
            localVarRequestOptions.data =  needsSerialization
                ? JSON.stringify(body !== undefined ? body : {})
                : (body || "");

            return {
                url: localVarUrlObj.pathname + localVarUrlObj.search + localVarUrlObj.hash,
                options: localVarRequestOptions,
            };
        },
        /**
         * 
         * @summary Lists SharedTaskDynamicZones
         * @param {string} [includes] Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names 
         * @param {string} [where] Filter on specific fields. Multiple conditions [.] separated Example: col_like_value.col2__val2
         * @param {string} [whereOr] Filter on specific fields (Chained ors). Multiple conditions [.] separated Example: col_like_value.col2__val2
         * @param {string} [groupBy] Group by field. Multiple conditions [.] separated Example: field1.field2
         * @param {string} [limit] Rows to limit in response (Default: 10,000)
         * @param {string} [orderBy] Order by [field]
         * @param {string} [orderDirection] Order by field direction
         * @param {string} [select] Column names [.] separated to fetch specific fields in response
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        listSharedTaskDynamicZones: async (includes?: string, where?: string, whereOr?: string, groupBy?: string, limit?: string, orderBy?: string, orderDirection?: string, select?: string, options: any = {}): Promise<RequestArgs> => {
            const localVarPath = `/shared_task_dynamic_zones`;
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, 'https://example.com');
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'GET', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;

            if (includes !== undefined) {
                localVarQueryParameter['includes'] = includes;
            }

            if (where !== undefined) {
                localVarQueryParameter['where'] = where;
            }

            if (whereOr !== undefined) {
                localVarQueryParameter['whereOr'] = whereOr;
            }

            if (groupBy !== undefined) {
                localVarQueryParameter['groupBy'] = groupBy;
            }

            if (limit !== undefined) {
                localVarQueryParameter['limit'] = limit;
            }

            if (orderBy !== undefined) {
                localVarQueryParameter['orderBy'] = orderBy;
            }

            if (orderDirection !== undefined) {
                localVarQueryParameter['orderDirection'] = orderDirection;
            }

            if (select !== undefined) {
                localVarQueryParameter['select'] = select;
            }


    
            const queryParameters = new URLSearchParams(localVarUrlObj.search);
            for (const key in localVarQueryParameter) {
                queryParameters.set(key, localVarQueryParameter[key]);
            }
            for (const key in options.query) {
                queryParameters.set(key, options.query[key]);
            }
            localVarUrlObj.search = (new URLSearchParams(queryParameters)).toString();
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};

            return {
                url: localVarUrlObj.pathname + localVarUrlObj.search + localVarUrlObj.hash,
                options: localVarRequestOptions,
            };
        },
        /**
         * 
         * @summary Updates SharedTaskDynamicZone
         * @param {number} id Id
         * @param {ModelsSharedTaskDynamicZone} sharedTaskDynamicZone SharedTaskDynamicZone
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        updateSharedTaskDynamicZone: async (id: number, sharedTaskDynamicZone: ModelsSharedTaskDynamicZone, options: any = {}): Promise<RequestArgs> => {
            // verify required parameter 'id' is not null or undefined
            if (id === null || id === undefined) {
                throw new RequiredError('id','Required parameter id was null or undefined when calling updateSharedTaskDynamicZone.');
            }
            // verify required parameter 'sharedTaskDynamicZone' is not null or undefined
            if (sharedTaskDynamicZone === null || sharedTaskDynamicZone === undefined) {
                throw new RequiredError('sharedTaskDynamicZone','Required parameter sharedTaskDynamicZone was null or undefined when calling updateSharedTaskDynamicZone.');
            }
            const localVarPath = `/shared_task_dynamic_zone/{id}`
                .replace(`{${"id"}}`, encodeURIComponent(String(id)));
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, 'https://example.com');
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'PATCH', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;


    
            localVarHeaderParameter['Content-Type'] = 'application/json';

            const queryParameters = new URLSearchParams(localVarUrlObj.search);
            for (const key in localVarQueryParameter) {
                queryParameters.set(key, localVarQueryParameter[key]);
            }
            for (const key in options.query) {
                queryParameters.set(key, options.query[key]);
            }
            localVarUrlObj.search = (new URLSearchParams(queryParameters)).toString();
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};
            const nonString = typeof sharedTaskDynamicZone !== 'string';
            const needsSerialization = nonString && configuration && configuration.isJsonMime
                ? configuration.isJsonMime(localVarRequestOptions.headers['Content-Type'])
                : nonString;
            localVarRequestOptions.data =  needsSerialization
                ? JSON.stringify(sharedTaskDynamicZone !== undefined ? sharedTaskDynamicZone : {})
                : (sharedTaskDynamicZone || "");

            return {
                url: localVarUrlObj.pathname + localVarUrlObj.search + localVarUrlObj.hash,
                options: localVarRequestOptions,
            };
        },
    }
};

/**
 * SharedTaskDynamicZoneApi - functional programming interface
 * @export
 */
export const SharedTaskDynamicZoneApiFp = function(configuration?: Configuration) {
    return {
        /**
         * 
         * @summary Creates SharedTaskDynamicZone
         * @param {ModelsSharedTaskDynamicZone} sharedTaskDynamicZone SharedTaskDynamicZone
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async createSharedTaskDynamicZone(sharedTaskDynamicZone: ModelsSharedTaskDynamicZone, options?: any): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<Array<ModelsSharedTaskDynamicZone>>> {
            const localVarAxiosArgs = await SharedTaskDynamicZoneApiAxiosParamCreator(configuration).createSharedTaskDynamicZone(sharedTaskDynamicZone, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs = {...localVarAxiosArgs.options, url: (configuration?.basePath || basePath) + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * 
         * @summary Deletes SharedTaskDynamicZone
         * @param {number} id sharedTaskId
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async deleteSharedTaskDynamicZone(id: number, options?: any): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<string>> {
            const localVarAxiosArgs = await SharedTaskDynamicZoneApiAxiosParamCreator(configuration).deleteSharedTaskDynamicZone(id, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs = {...localVarAxiosArgs.options, url: (configuration?.basePath || basePath) + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * 
         * @summary Gets SharedTaskDynamicZone
         * @param {number} id Id
         * @param {string} [includes] Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names 
         * @param {string} [select] Column names [.] separated to fetch specific fields in response
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getSharedTaskDynamicZone(id: number, includes?: string, select?: string, options?: any): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<Array<ModelsSharedTaskDynamicZone>>> {
            const localVarAxiosArgs = await SharedTaskDynamicZoneApiAxiosParamCreator(configuration).getSharedTaskDynamicZone(id, includes, select, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs = {...localVarAxiosArgs.options, url: (configuration?.basePath || basePath) + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * 
         * @summary Gets SharedTaskDynamicZones in bulk
         * @param {CrudcontrollersBulkFetchByIdsGetRequest} body body
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getSharedTaskDynamicZonesBulk(body: CrudcontrollersBulkFetchByIdsGetRequest, options?: any): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<Array<ModelsSharedTaskDynamicZone>>> {
            const localVarAxiosArgs = await SharedTaskDynamicZoneApiAxiosParamCreator(configuration).getSharedTaskDynamicZonesBulk(body, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs = {...localVarAxiosArgs.options, url: (configuration?.basePath || basePath) + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * 
         * @summary Lists SharedTaskDynamicZones
         * @param {string} [includes] Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names 
         * @param {string} [where] Filter on specific fields. Multiple conditions [.] separated Example: col_like_value.col2__val2
         * @param {string} [whereOr] Filter on specific fields (Chained ors). Multiple conditions [.] separated Example: col_like_value.col2__val2
         * @param {string} [groupBy] Group by field. Multiple conditions [.] separated Example: field1.field2
         * @param {string} [limit] Rows to limit in response (Default: 10,000)
         * @param {string} [orderBy] Order by [field]
         * @param {string} [orderDirection] Order by field direction
         * @param {string} [select] Column names [.] separated to fetch specific fields in response
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async listSharedTaskDynamicZones(includes?: string, where?: string, whereOr?: string, groupBy?: string, limit?: string, orderBy?: string, orderDirection?: string, select?: string, options?: any): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<Array<ModelsSharedTaskDynamicZone>>> {
            const localVarAxiosArgs = await SharedTaskDynamicZoneApiAxiosParamCreator(configuration).listSharedTaskDynamicZones(includes, where, whereOr, groupBy, limit, orderBy, orderDirection, select, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs = {...localVarAxiosArgs.options, url: (configuration?.basePath || basePath) + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * 
         * @summary Updates SharedTaskDynamicZone
         * @param {number} id Id
         * @param {ModelsSharedTaskDynamicZone} sharedTaskDynamicZone SharedTaskDynamicZone
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async updateSharedTaskDynamicZone(id: number, sharedTaskDynamicZone: ModelsSharedTaskDynamicZone, options?: any): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<Array<ModelsSharedTaskDynamicZone>>> {
            const localVarAxiosArgs = await SharedTaskDynamicZoneApiAxiosParamCreator(configuration).updateSharedTaskDynamicZone(id, sharedTaskDynamicZone, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs = {...localVarAxiosArgs.options, url: (configuration?.basePath || basePath) + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
    }
};

/**
 * SharedTaskDynamicZoneApi - factory interface
 * @export
 */
export const SharedTaskDynamicZoneApiFactory = function (configuration?: Configuration, basePath?: string, axios?: AxiosInstance) {
    return {
        /**
         * 
         * @summary Creates SharedTaskDynamicZone
         * @param {ModelsSharedTaskDynamicZone} sharedTaskDynamicZone SharedTaskDynamicZone
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        createSharedTaskDynamicZone(sharedTaskDynamicZone: ModelsSharedTaskDynamicZone, options?: any): AxiosPromise<Array<ModelsSharedTaskDynamicZone>> {
            return SharedTaskDynamicZoneApiFp(configuration).createSharedTaskDynamicZone(sharedTaskDynamicZone, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Deletes SharedTaskDynamicZone
         * @param {number} id sharedTaskId
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        deleteSharedTaskDynamicZone(id: number, options?: any): AxiosPromise<string> {
            return SharedTaskDynamicZoneApiFp(configuration).deleteSharedTaskDynamicZone(id, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Gets SharedTaskDynamicZone
         * @param {number} id Id
         * @param {string} [includes] Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names 
         * @param {string} [select] Column names [.] separated to fetch specific fields in response
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getSharedTaskDynamicZone(id: number, includes?: string, select?: string, options?: any): AxiosPromise<Array<ModelsSharedTaskDynamicZone>> {
            return SharedTaskDynamicZoneApiFp(configuration).getSharedTaskDynamicZone(id, includes, select, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Gets SharedTaskDynamicZones in bulk
         * @param {CrudcontrollersBulkFetchByIdsGetRequest} body body
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getSharedTaskDynamicZonesBulk(body: CrudcontrollersBulkFetchByIdsGetRequest, options?: any): AxiosPromise<Array<ModelsSharedTaskDynamicZone>> {
            return SharedTaskDynamicZoneApiFp(configuration).getSharedTaskDynamicZonesBulk(body, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Lists SharedTaskDynamicZones
         * @param {string} [includes] Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names 
         * @param {string} [where] Filter on specific fields. Multiple conditions [.] separated Example: col_like_value.col2__val2
         * @param {string} [whereOr] Filter on specific fields (Chained ors). Multiple conditions [.] separated Example: col_like_value.col2__val2
         * @param {string} [groupBy] Group by field. Multiple conditions [.] separated Example: field1.field2
         * @param {string} [limit] Rows to limit in response (Default: 10,000)
         * @param {string} [orderBy] Order by [field]
         * @param {string} [orderDirection] Order by field direction
         * @param {string} [select] Column names [.] separated to fetch specific fields in response
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        listSharedTaskDynamicZones(includes?: string, where?: string, whereOr?: string, groupBy?: string, limit?: string, orderBy?: string, orderDirection?: string, select?: string, options?: any): AxiosPromise<Array<ModelsSharedTaskDynamicZone>> {
            return SharedTaskDynamicZoneApiFp(configuration).listSharedTaskDynamicZones(includes, where, whereOr, groupBy, limit, orderBy, orderDirection, select, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Updates SharedTaskDynamicZone
         * @param {number} id Id
         * @param {ModelsSharedTaskDynamicZone} sharedTaskDynamicZone SharedTaskDynamicZone
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        updateSharedTaskDynamicZone(id: number, sharedTaskDynamicZone: ModelsSharedTaskDynamicZone, options?: any): AxiosPromise<Array<ModelsSharedTaskDynamicZone>> {
            return SharedTaskDynamicZoneApiFp(configuration).updateSharedTaskDynamicZone(id, sharedTaskDynamicZone, options).then((request) => request(axios, basePath));
        },
    };
};

/**
 * Request parameters for createSharedTaskDynamicZone operation in SharedTaskDynamicZoneApi.
 * @export
 * @interface SharedTaskDynamicZoneApiCreateSharedTaskDynamicZoneRequest
 */
export interface SharedTaskDynamicZoneApiCreateSharedTaskDynamicZoneRequest {
    /**
     * SharedTaskDynamicZone
     * @type {ModelsSharedTaskDynamicZone}
     * @memberof SharedTaskDynamicZoneApiCreateSharedTaskDynamicZone
     */
    readonly sharedTaskDynamicZone: ModelsSharedTaskDynamicZone
}

/**
 * Request parameters for deleteSharedTaskDynamicZone operation in SharedTaskDynamicZoneApi.
 * @export
 * @interface SharedTaskDynamicZoneApiDeleteSharedTaskDynamicZoneRequest
 */
export interface SharedTaskDynamicZoneApiDeleteSharedTaskDynamicZoneRequest {
    /**
     * sharedTaskId
     * @type {number}
     * @memberof SharedTaskDynamicZoneApiDeleteSharedTaskDynamicZone
     */
    readonly id: number
}

/**
 * Request parameters for getSharedTaskDynamicZone operation in SharedTaskDynamicZoneApi.
 * @export
 * @interface SharedTaskDynamicZoneApiGetSharedTaskDynamicZoneRequest
 */
export interface SharedTaskDynamicZoneApiGetSharedTaskDynamicZoneRequest {
    /**
     * Id
     * @type {number}
     * @memberof SharedTaskDynamicZoneApiGetSharedTaskDynamicZone
     */
    readonly id: number

    /**
     * Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names 
     * @type {string}
     * @memberof SharedTaskDynamicZoneApiGetSharedTaskDynamicZone
     */
    readonly includes?: string

    /**
     * Column names [.] separated to fetch specific fields in response
     * @type {string}
     * @memberof SharedTaskDynamicZoneApiGetSharedTaskDynamicZone
     */
    readonly select?: string
}

/**
 * Request parameters for getSharedTaskDynamicZonesBulk operation in SharedTaskDynamicZoneApi.
 * @export
 * @interface SharedTaskDynamicZoneApiGetSharedTaskDynamicZonesBulkRequest
 */
export interface SharedTaskDynamicZoneApiGetSharedTaskDynamicZonesBulkRequest {
    /**
     * body
     * @type {CrudcontrollersBulkFetchByIdsGetRequest}
     * @memberof SharedTaskDynamicZoneApiGetSharedTaskDynamicZonesBulk
     */
    readonly body: CrudcontrollersBulkFetchByIdsGetRequest
}

/**
 * Request parameters for listSharedTaskDynamicZones operation in SharedTaskDynamicZoneApi.
 * @export
 * @interface SharedTaskDynamicZoneApiListSharedTaskDynamicZonesRequest
 */
export interface SharedTaskDynamicZoneApiListSharedTaskDynamicZonesRequest {
    /**
     * Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names 
     * @type {string}
     * @memberof SharedTaskDynamicZoneApiListSharedTaskDynamicZones
     */
    readonly includes?: string

    /**
     * Filter on specific fields. Multiple conditions [.] separated Example: col_like_value.col2__val2
     * @type {string}
     * @memberof SharedTaskDynamicZoneApiListSharedTaskDynamicZones
     */
    readonly where?: string

    /**
     * Filter on specific fields (Chained ors). Multiple conditions [.] separated Example: col_like_value.col2__val2
     * @type {string}
     * @memberof SharedTaskDynamicZoneApiListSharedTaskDynamicZones
     */
    readonly whereOr?: string

    /**
     * Group by field. Multiple conditions [.] separated Example: field1.field2
     * @type {string}
     * @memberof SharedTaskDynamicZoneApiListSharedTaskDynamicZones
     */
    readonly groupBy?: string

    /**
     * Rows to limit in response (Default: 10,000)
     * @type {string}
     * @memberof SharedTaskDynamicZoneApiListSharedTaskDynamicZones
     */
    readonly limit?: string

    /**
     * Order by [field]
     * @type {string}
     * @memberof SharedTaskDynamicZoneApiListSharedTaskDynamicZones
     */
    readonly orderBy?: string

    /**
     * Order by field direction
     * @type {string}
     * @memberof SharedTaskDynamicZoneApiListSharedTaskDynamicZones
     */
    readonly orderDirection?: string

    /**
     * Column names [.] separated to fetch specific fields in response
     * @type {string}
     * @memberof SharedTaskDynamicZoneApiListSharedTaskDynamicZones
     */
    readonly select?: string
}

/**
 * Request parameters for updateSharedTaskDynamicZone operation in SharedTaskDynamicZoneApi.
 * @export
 * @interface SharedTaskDynamicZoneApiUpdateSharedTaskDynamicZoneRequest
 */
export interface SharedTaskDynamicZoneApiUpdateSharedTaskDynamicZoneRequest {
    /**
     * Id
     * @type {number}
     * @memberof SharedTaskDynamicZoneApiUpdateSharedTaskDynamicZone
     */
    readonly id: number

    /**
     * SharedTaskDynamicZone
     * @type {ModelsSharedTaskDynamicZone}
     * @memberof SharedTaskDynamicZoneApiUpdateSharedTaskDynamicZone
     */
    readonly sharedTaskDynamicZone: ModelsSharedTaskDynamicZone
}

/**
 * SharedTaskDynamicZoneApi - object-oriented interface
 * @export
 * @class SharedTaskDynamicZoneApi
 * @extends {BaseAPI}
 */
export class SharedTaskDynamicZoneApi extends BaseAPI {
    /**
     * 
     * @summary Creates SharedTaskDynamicZone
     * @param {SharedTaskDynamicZoneApiCreateSharedTaskDynamicZoneRequest} requestParameters Request parameters.
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof SharedTaskDynamicZoneApi
     */
    public createSharedTaskDynamicZone(requestParameters: SharedTaskDynamicZoneApiCreateSharedTaskDynamicZoneRequest, options?: any) {
        return SharedTaskDynamicZoneApiFp(this.configuration).createSharedTaskDynamicZone(requestParameters.sharedTaskDynamicZone, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Deletes SharedTaskDynamicZone
     * @param {SharedTaskDynamicZoneApiDeleteSharedTaskDynamicZoneRequest} requestParameters Request parameters.
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof SharedTaskDynamicZoneApi
     */
    public deleteSharedTaskDynamicZone(requestParameters: SharedTaskDynamicZoneApiDeleteSharedTaskDynamicZoneRequest, options?: any) {
        return SharedTaskDynamicZoneApiFp(this.configuration).deleteSharedTaskDynamicZone(requestParameters.id, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Gets SharedTaskDynamicZone
     * @param {SharedTaskDynamicZoneApiGetSharedTaskDynamicZoneRequest} requestParameters Request parameters.
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof SharedTaskDynamicZoneApi
     */
    public getSharedTaskDynamicZone(requestParameters: SharedTaskDynamicZoneApiGetSharedTaskDynamicZoneRequest, options?: any) {
        return SharedTaskDynamicZoneApiFp(this.configuration).getSharedTaskDynamicZone(requestParameters.id, requestParameters.includes, requestParameters.select, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Gets SharedTaskDynamicZones in bulk
     * @param {SharedTaskDynamicZoneApiGetSharedTaskDynamicZonesBulkRequest} requestParameters Request parameters.
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof SharedTaskDynamicZoneApi
     */
    public getSharedTaskDynamicZonesBulk(requestParameters: SharedTaskDynamicZoneApiGetSharedTaskDynamicZonesBulkRequest, options?: any) {
        return SharedTaskDynamicZoneApiFp(this.configuration).getSharedTaskDynamicZonesBulk(requestParameters.body, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Lists SharedTaskDynamicZones
     * @param {SharedTaskDynamicZoneApiListSharedTaskDynamicZonesRequest} requestParameters Request parameters.
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof SharedTaskDynamicZoneApi
     */
    public listSharedTaskDynamicZones(requestParameters: SharedTaskDynamicZoneApiListSharedTaskDynamicZonesRequest = {}, options?: any) {
        return SharedTaskDynamicZoneApiFp(this.configuration).listSharedTaskDynamicZones(requestParameters.includes, requestParameters.where, requestParameters.whereOr, requestParameters.groupBy, requestParameters.limit, requestParameters.orderBy, requestParameters.orderDirection, requestParameters.select, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Updates SharedTaskDynamicZone
     * @param {SharedTaskDynamicZoneApiUpdateSharedTaskDynamicZoneRequest} requestParameters Request parameters.
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof SharedTaskDynamicZoneApi
     */
    public updateSharedTaskDynamicZone(requestParameters: SharedTaskDynamicZoneApiUpdateSharedTaskDynamicZoneRequest, options?: any) {
        return SharedTaskDynamicZoneApiFp(this.configuration).updateSharedTaskDynamicZone(requestParameters.id, requestParameters.sharedTaskDynamicZone, options).then((request) => request(this.axios, this.basePath));
    }
}