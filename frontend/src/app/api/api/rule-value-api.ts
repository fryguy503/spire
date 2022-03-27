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
import { ModelsRuleValue } from '../models';
/**
 * RuleValueApi - axios parameter creator
 * @export
 */
export const RuleValueApiAxiosParamCreator = function (configuration?: Configuration) {
    return {
        /**
         * 
         * @summary Creates RuleValue
         * @param {ModelsRuleValue} ruleValue RuleValue
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        createRuleValue: async (ruleValue: ModelsRuleValue, options: any = {}): Promise<RequestArgs> => {
            // verify required parameter 'ruleValue' is not null or undefined
            if (ruleValue === null || ruleValue === undefined) {
                throw new RequiredError('ruleValue','Required parameter ruleValue was null or undefined when calling createRuleValue.');
            }
            const localVarPath = `/rule_value`;
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
            const nonString = typeof ruleValue !== 'string';
            const needsSerialization = nonString && configuration && configuration.isJsonMime
                ? configuration.isJsonMime(localVarRequestOptions.headers['Content-Type'])
                : nonString;
            localVarRequestOptions.data =  needsSerialization
                ? JSON.stringify(ruleValue !== undefined ? ruleValue : {})
                : (ruleValue || "");

            return {
                url: localVarUrlObj.pathname + localVarUrlObj.search + localVarUrlObj.hash,
                options: localVarRequestOptions,
            };
        },
        /**
         * 
         * @summary Deletes RuleValue
         * @param {number} id rulesetId
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        deleteRuleValue: async (id: number, options: any = {}): Promise<RequestArgs> => {
            // verify required parameter 'id' is not null or undefined
            if (id === null || id === undefined) {
                throw new RequiredError('id','Required parameter id was null or undefined when calling deleteRuleValue.');
            }
            const localVarPath = `/rule_value/{id}`
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
         * @summary Gets RuleValue
         * @param {number} id Id
         * @param {string} [includes] Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names 
         * @param {string} [select] Column names [.] separated to fetch specific fields in response
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getRuleValue: async (id: number, includes?: string, select?: string, options: any = {}): Promise<RequestArgs> => {
            // verify required parameter 'id' is not null or undefined
            if (id === null || id === undefined) {
                throw new RequiredError('id','Required parameter id was null or undefined when calling getRuleValue.');
            }
            const localVarPath = `/rule_value/{id}`
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
         * @summary Gets RuleValues in bulk
         * @param {CrudcontrollersBulkFetchByIdsGetRequest} body body
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getRuleValuesBulk: async (body: CrudcontrollersBulkFetchByIdsGetRequest, options: any = {}): Promise<RequestArgs> => {
            // verify required parameter 'body' is not null or undefined
            if (body === null || body === undefined) {
                throw new RequiredError('body','Required parameter body was null or undefined when calling getRuleValuesBulk.');
            }
            const localVarPath = `/rule_values/bulk`;
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
         * @summary Lists RuleValues
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
        listRuleValues: async (includes?: string, where?: string, whereOr?: string, groupBy?: string, limit?: string, orderBy?: string, orderDirection?: string, select?: string, options: any = {}): Promise<RequestArgs> => {
            const localVarPath = `/rule_values`;
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
         * @summary Updates RuleValue
         * @param {number} id Id
         * @param {ModelsRuleValue} ruleValue RuleValue
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        updateRuleValue: async (id: number, ruleValue: ModelsRuleValue, options: any = {}): Promise<RequestArgs> => {
            // verify required parameter 'id' is not null or undefined
            if (id === null || id === undefined) {
                throw new RequiredError('id','Required parameter id was null or undefined when calling updateRuleValue.');
            }
            // verify required parameter 'ruleValue' is not null or undefined
            if (ruleValue === null || ruleValue === undefined) {
                throw new RequiredError('ruleValue','Required parameter ruleValue was null or undefined when calling updateRuleValue.');
            }
            const localVarPath = `/rule_value/{id}`
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
            const nonString = typeof ruleValue !== 'string';
            const needsSerialization = nonString && configuration && configuration.isJsonMime
                ? configuration.isJsonMime(localVarRequestOptions.headers['Content-Type'])
                : nonString;
            localVarRequestOptions.data =  needsSerialization
                ? JSON.stringify(ruleValue !== undefined ? ruleValue : {})
                : (ruleValue || "");

            return {
                url: localVarUrlObj.pathname + localVarUrlObj.search + localVarUrlObj.hash,
                options: localVarRequestOptions,
            };
        },
    }
};

/**
 * RuleValueApi - functional programming interface
 * @export
 */
export const RuleValueApiFp = function(configuration?: Configuration) {
    return {
        /**
         * 
         * @summary Creates RuleValue
         * @param {ModelsRuleValue} ruleValue RuleValue
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async createRuleValue(ruleValue: ModelsRuleValue, options?: any): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<Array<ModelsRuleValue>>> {
            const localVarAxiosArgs = await RuleValueApiAxiosParamCreator(configuration).createRuleValue(ruleValue, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs = {...localVarAxiosArgs.options, url: (configuration?.basePath || basePath) + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * 
         * @summary Deletes RuleValue
         * @param {number} id rulesetId
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async deleteRuleValue(id: number, options?: any): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<string>> {
            const localVarAxiosArgs = await RuleValueApiAxiosParamCreator(configuration).deleteRuleValue(id, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs = {...localVarAxiosArgs.options, url: (configuration?.basePath || basePath) + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * 
         * @summary Gets RuleValue
         * @param {number} id Id
         * @param {string} [includes] Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names 
         * @param {string} [select] Column names [.] separated to fetch specific fields in response
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getRuleValue(id: number, includes?: string, select?: string, options?: any): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<Array<ModelsRuleValue>>> {
            const localVarAxiosArgs = await RuleValueApiAxiosParamCreator(configuration).getRuleValue(id, includes, select, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs = {...localVarAxiosArgs.options, url: (configuration?.basePath || basePath) + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * 
         * @summary Gets RuleValues in bulk
         * @param {CrudcontrollersBulkFetchByIdsGetRequest} body body
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getRuleValuesBulk(body: CrudcontrollersBulkFetchByIdsGetRequest, options?: any): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<Array<ModelsRuleValue>>> {
            const localVarAxiosArgs = await RuleValueApiAxiosParamCreator(configuration).getRuleValuesBulk(body, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs = {...localVarAxiosArgs.options, url: (configuration?.basePath || basePath) + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * 
         * @summary Lists RuleValues
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
        async listRuleValues(includes?: string, where?: string, whereOr?: string, groupBy?: string, limit?: string, orderBy?: string, orderDirection?: string, select?: string, options?: any): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<Array<ModelsRuleValue>>> {
            const localVarAxiosArgs = await RuleValueApiAxiosParamCreator(configuration).listRuleValues(includes, where, whereOr, groupBy, limit, orderBy, orderDirection, select, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs = {...localVarAxiosArgs.options, url: (configuration?.basePath || basePath) + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * 
         * @summary Updates RuleValue
         * @param {number} id Id
         * @param {ModelsRuleValue} ruleValue RuleValue
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async updateRuleValue(id: number, ruleValue: ModelsRuleValue, options?: any): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<Array<ModelsRuleValue>>> {
            const localVarAxiosArgs = await RuleValueApiAxiosParamCreator(configuration).updateRuleValue(id, ruleValue, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs = {...localVarAxiosArgs.options, url: (configuration?.basePath || basePath) + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
    }
};

/**
 * RuleValueApi - factory interface
 * @export
 */
export const RuleValueApiFactory = function (configuration?: Configuration, basePath?: string, axios?: AxiosInstance) {
    return {
        /**
         * 
         * @summary Creates RuleValue
         * @param {ModelsRuleValue} ruleValue RuleValue
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        createRuleValue(ruleValue: ModelsRuleValue, options?: any): AxiosPromise<Array<ModelsRuleValue>> {
            return RuleValueApiFp(configuration).createRuleValue(ruleValue, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Deletes RuleValue
         * @param {number} id rulesetId
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        deleteRuleValue(id: number, options?: any): AxiosPromise<string> {
            return RuleValueApiFp(configuration).deleteRuleValue(id, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Gets RuleValue
         * @param {number} id Id
         * @param {string} [includes] Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names 
         * @param {string} [select] Column names [.] separated to fetch specific fields in response
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getRuleValue(id: number, includes?: string, select?: string, options?: any): AxiosPromise<Array<ModelsRuleValue>> {
            return RuleValueApiFp(configuration).getRuleValue(id, includes, select, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Gets RuleValues in bulk
         * @param {CrudcontrollersBulkFetchByIdsGetRequest} body body
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getRuleValuesBulk(body: CrudcontrollersBulkFetchByIdsGetRequest, options?: any): AxiosPromise<Array<ModelsRuleValue>> {
            return RuleValueApiFp(configuration).getRuleValuesBulk(body, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Lists RuleValues
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
        listRuleValues(includes?: string, where?: string, whereOr?: string, groupBy?: string, limit?: string, orderBy?: string, orderDirection?: string, select?: string, options?: any): AxiosPromise<Array<ModelsRuleValue>> {
            return RuleValueApiFp(configuration).listRuleValues(includes, where, whereOr, groupBy, limit, orderBy, orderDirection, select, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Updates RuleValue
         * @param {number} id Id
         * @param {ModelsRuleValue} ruleValue RuleValue
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        updateRuleValue(id: number, ruleValue: ModelsRuleValue, options?: any): AxiosPromise<Array<ModelsRuleValue>> {
            return RuleValueApiFp(configuration).updateRuleValue(id, ruleValue, options).then((request) => request(axios, basePath));
        },
    };
};

/**
 * Request parameters for createRuleValue operation in RuleValueApi.
 * @export
 * @interface RuleValueApiCreateRuleValueRequest
 */
export interface RuleValueApiCreateRuleValueRequest {
    /**
     * RuleValue
     * @type {ModelsRuleValue}
     * @memberof RuleValueApiCreateRuleValue
     */
    readonly ruleValue: ModelsRuleValue
}

/**
 * Request parameters for deleteRuleValue operation in RuleValueApi.
 * @export
 * @interface RuleValueApiDeleteRuleValueRequest
 */
export interface RuleValueApiDeleteRuleValueRequest {
    /**
     * rulesetId
     * @type {number}
     * @memberof RuleValueApiDeleteRuleValue
     */
    readonly id: number
}

/**
 * Request parameters for getRuleValue operation in RuleValueApi.
 * @export
 * @interface RuleValueApiGetRuleValueRequest
 */
export interface RuleValueApiGetRuleValueRequest {
    /**
     * Id
     * @type {number}
     * @memberof RuleValueApiGetRuleValue
     */
    readonly id: number

    /**
     * Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names 
     * @type {string}
     * @memberof RuleValueApiGetRuleValue
     */
    readonly includes?: string

    /**
     * Column names [.] separated to fetch specific fields in response
     * @type {string}
     * @memberof RuleValueApiGetRuleValue
     */
    readonly select?: string
}

/**
 * Request parameters for getRuleValuesBulk operation in RuleValueApi.
 * @export
 * @interface RuleValueApiGetRuleValuesBulkRequest
 */
export interface RuleValueApiGetRuleValuesBulkRequest {
    /**
     * body
     * @type {CrudcontrollersBulkFetchByIdsGetRequest}
     * @memberof RuleValueApiGetRuleValuesBulk
     */
    readonly body: CrudcontrollersBulkFetchByIdsGetRequest
}

/**
 * Request parameters for listRuleValues operation in RuleValueApi.
 * @export
 * @interface RuleValueApiListRuleValuesRequest
 */
export interface RuleValueApiListRuleValuesRequest {
    /**
     * Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names 
     * @type {string}
     * @memberof RuleValueApiListRuleValues
     */
    readonly includes?: string

    /**
     * Filter on specific fields. Multiple conditions [.] separated Example: col_like_value.col2__val2
     * @type {string}
     * @memberof RuleValueApiListRuleValues
     */
    readonly where?: string

    /**
     * Filter on specific fields (Chained ors). Multiple conditions [.] separated Example: col_like_value.col2__val2
     * @type {string}
     * @memberof RuleValueApiListRuleValues
     */
    readonly whereOr?: string

    /**
     * Group by field. Multiple conditions [.] separated Example: field1.field2
     * @type {string}
     * @memberof RuleValueApiListRuleValues
     */
    readonly groupBy?: string

    /**
     * Rows to limit in response (Default: 10,000)
     * @type {string}
     * @memberof RuleValueApiListRuleValues
     */
    readonly limit?: string

    /**
     * Order by [field]
     * @type {string}
     * @memberof RuleValueApiListRuleValues
     */
    readonly orderBy?: string

    /**
     * Order by field direction
     * @type {string}
     * @memberof RuleValueApiListRuleValues
     */
    readonly orderDirection?: string

    /**
     * Column names [.] separated to fetch specific fields in response
     * @type {string}
     * @memberof RuleValueApiListRuleValues
     */
    readonly select?: string
}

/**
 * Request parameters for updateRuleValue operation in RuleValueApi.
 * @export
 * @interface RuleValueApiUpdateRuleValueRequest
 */
export interface RuleValueApiUpdateRuleValueRequest {
    /**
     * Id
     * @type {number}
     * @memberof RuleValueApiUpdateRuleValue
     */
    readonly id: number

    /**
     * RuleValue
     * @type {ModelsRuleValue}
     * @memberof RuleValueApiUpdateRuleValue
     */
    readonly ruleValue: ModelsRuleValue
}

/**
 * RuleValueApi - object-oriented interface
 * @export
 * @class RuleValueApi
 * @extends {BaseAPI}
 */
export class RuleValueApi extends BaseAPI {
    /**
     * 
     * @summary Creates RuleValue
     * @param {RuleValueApiCreateRuleValueRequest} requestParameters Request parameters.
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof RuleValueApi
     */
    public createRuleValue(requestParameters: RuleValueApiCreateRuleValueRequest, options?: any) {
        return RuleValueApiFp(this.configuration).createRuleValue(requestParameters.ruleValue, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Deletes RuleValue
     * @param {RuleValueApiDeleteRuleValueRequest} requestParameters Request parameters.
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof RuleValueApi
     */
    public deleteRuleValue(requestParameters: RuleValueApiDeleteRuleValueRequest, options?: any) {
        return RuleValueApiFp(this.configuration).deleteRuleValue(requestParameters.id, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Gets RuleValue
     * @param {RuleValueApiGetRuleValueRequest} requestParameters Request parameters.
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof RuleValueApi
     */
    public getRuleValue(requestParameters: RuleValueApiGetRuleValueRequest, options?: any) {
        return RuleValueApiFp(this.configuration).getRuleValue(requestParameters.id, requestParameters.includes, requestParameters.select, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Gets RuleValues in bulk
     * @param {RuleValueApiGetRuleValuesBulkRequest} requestParameters Request parameters.
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof RuleValueApi
     */
    public getRuleValuesBulk(requestParameters: RuleValueApiGetRuleValuesBulkRequest, options?: any) {
        return RuleValueApiFp(this.configuration).getRuleValuesBulk(requestParameters.body, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Lists RuleValues
     * @param {RuleValueApiListRuleValuesRequest} requestParameters Request parameters.
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof RuleValueApi
     */
    public listRuleValues(requestParameters: RuleValueApiListRuleValuesRequest = {}, options?: any) {
        return RuleValueApiFp(this.configuration).listRuleValues(requestParameters.includes, requestParameters.where, requestParameters.whereOr, requestParameters.groupBy, requestParameters.limit, requestParameters.orderBy, requestParameters.orderDirection, requestParameters.select, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Updates RuleValue
     * @param {RuleValueApiUpdateRuleValueRequest} requestParameters Request parameters.
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof RuleValueApi
     */
    public updateRuleValue(requestParameters: RuleValueApiUpdateRuleValueRequest, options?: any) {
        return RuleValueApiFp(this.configuration).updateRuleValue(requestParameters.id, requestParameters.ruleValue, options).then((request) => request(this.axios, this.basePath));
    }
}