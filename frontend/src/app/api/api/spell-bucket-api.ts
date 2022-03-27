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
import { ModelsSpellBucket } from '../models';
/**
 * SpellBucketApi - axios parameter creator
 * @export
 */
export const SpellBucketApiAxiosParamCreator = function (configuration?: Configuration) {
    return {
        /**
         * 
         * @summary Creates SpellBucket
         * @param {ModelsSpellBucket} spellBucket SpellBucket
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        createSpellBucket: async (spellBucket: ModelsSpellBucket, options: any = {}): Promise<RequestArgs> => {
            // verify required parameter 'spellBucket' is not null or undefined
            if (spellBucket === null || spellBucket === undefined) {
                throw new RequiredError('spellBucket','Required parameter spellBucket was null or undefined when calling createSpellBucket.');
            }
            const localVarPath = `/spell_bucket`;
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
            const nonString = typeof spellBucket !== 'string';
            const needsSerialization = nonString && configuration && configuration.isJsonMime
                ? configuration.isJsonMime(localVarRequestOptions.headers['Content-Type'])
                : nonString;
            localVarRequestOptions.data =  needsSerialization
                ? JSON.stringify(spellBucket !== undefined ? spellBucket : {})
                : (spellBucket || "");

            return {
                url: localVarUrlObj.pathname + localVarUrlObj.search + localVarUrlObj.hash,
                options: localVarRequestOptions,
            };
        },
        /**
         * 
         * @summary Deletes SpellBucket
         * @param {number} id spellid
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        deleteSpellBucket: async (id: number, options: any = {}): Promise<RequestArgs> => {
            // verify required parameter 'id' is not null or undefined
            if (id === null || id === undefined) {
                throw new RequiredError('id','Required parameter id was null or undefined when calling deleteSpellBucket.');
            }
            const localVarPath = `/spell_bucket/{id}`
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
         * @summary Gets SpellBucket
         * @param {number} id Id
         * @param {string} [includes] Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names 
         * @param {string} [select] Column names [.] separated to fetch specific fields in response
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getSpellBucket: async (id: number, includes?: string, select?: string, options: any = {}): Promise<RequestArgs> => {
            // verify required parameter 'id' is not null or undefined
            if (id === null || id === undefined) {
                throw new RequiredError('id','Required parameter id was null or undefined when calling getSpellBucket.');
            }
            const localVarPath = `/spell_bucket/{id}`
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
         * @summary Gets SpellBuckets in bulk
         * @param {CrudcontrollersBulkFetchByIdsGetRequest} body body
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getSpellBucketsBulk: async (body: CrudcontrollersBulkFetchByIdsGetRequest, options: any = {}): Promise<RequestArgs> => {
            // verify required parameter 'body' is not null or undefined
            if (body === null || body === undefined) {
                throw new RequiredError('body','Required parameter body was null or undefined when calling getSpellBucketsBulk.');
            }
            const localVarPath = `/spell_buckets/bulk`;
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
         * @summary Lists SpellBuckets
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
        listSpellBuckets: async (includes?: string, where?: string, whereOr?: string, groupBy?: string, limit?: string, orderBy?: string, orderDirection?: string, select?: string, options: any = {}): Promise<RequestArgs> => {
            const localVarPath = `/spell_buckets`;
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
         * @summary Updates SpellBucket
         * @param {number} id Id
         * @param {ModelsSpellBucket} spellBucket SpellBucket
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        updateSpellBucket: async (id: number, spellBucket: ModelsSpellBucket, options: any = {}): Promise<RequestArgs> => {
            // verify required parameter 'id' is not null or undefined
            if (id === null || id === undefined) {
                throw new RequiredError('id','Required parameter id was null or undefined when calling updateSpellBucket.');
            }
            // verify required parameter 'spellBucket' is not null or undefined
            if (spellBucket === null || spellBucket === undefined) {
                throw new RequiredError('spellBucket','Required parameter spellBucket was null or undefined when calling updateSpellBucket.');
            }
            const localVarPath = `/spell_bucket/{id}`
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
            const nonString = typeof spellBucket !== 'string';
            const needsSerialization = nonString && configuration && configuration.isJsonMime
                ? configuration.isJsonMime(localVarRequestOptions.headers['Content-Type'])
                : nonString;
            localVarRequestOptions.data =  needsSerialization
                ? JSON.stringify(spellBucket !== undefined ? spellBucket : {})
                : (spellBucket || "");

            return {
                url: localVarUrlObj.pathname + localVarUrlObj.search + localVarUrlObj.hash,
                options: localVarRequestOptions,
            };
        },
    }
};

/**
 * SpellBucketApi - functional programming interface
 * @export
 */
export const SpellBucketApiFp = function(configuration?: Configuration) {
    return {
        /**
         * 
         * @summary Creates SpellBucket
         * @param {ModelsSpellBucket} spellBucket SpellBucket
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async createSpellBucket(spellBucket: ModelsSpellBucket, options?: any): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<Array<ModelsSpellBucket>>> {
            const localVarAxiosArgs = await SpellBucketApiAxiosParamCreator(configuration).createSpellBucket(spellBucket, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs = {...localVarAxiosArgs.options, url: (configuration?.basePath || basePath) + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * 
         * @summary Deletes SpellBucket
         * @param {number} id spellid
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async deleteSpellBucket(id: number, options?: any): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<string>> {
            const localVarAxiosArgs = await SpellBucketApiAxiosParamCreator(configuration).deleteSpellBucket(id, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs = {...localVarAxiosArgs.options, url: (configuration?.basePath || basePath) + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * 
         * @summary Gets SpellBucket
         * @param {number} id Id
         * @param {string} [includes] Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names 
         * @param {string} [select] Column names [.] separated to fetch specific fields in response
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getSpellBucket(id: number, includes?: string, select?: string, options?: any): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<Array<ModelsSpellBucket>>> {
            const localVarAxiosArgs = await SpellBucketApiAxiosParamCreator(configuration).getSpellBucket(id, includes, select, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs = {...localVarAxiosArgs.options, url: (configuration?.basePath || basePath) + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * 
         * @summary Gets SpellBuckets in bulk
         * @param {CrudcontrollersBulkFetchByIdsGetRequest} body body
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getSpellBucketsBulk(body: CrudcontrollersBulkFetchByIdsGetRequest, options?: any): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<Array<ModelsSpellBucket>>> {
            const localVarAxiosArgs = await SpellBucketApiAxiosParamCreator(configuration).getSpellBucketsBulk(body, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs = {...localVarAxiosArgs.options, url: (configuration?.basePath || basePath) + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * 
         * @summary Lists SpellBuckets
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
        async listSpellBuckets(includes?: string, where?: string, whereOr?: string, groupBy?: string, limit?: string, orderBy?: string, orderDirection?: string, select?: string, options?: any): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<Array<ModelsSpellBucket>>> {
            const localVarAxiosArgs = await SpellBucketApiAxiosParamCreator(configuration).listSpellBuckets(includes, where, whereOr, groupBy, limit, orderBy, orderDirection, select, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs = {...localVarAxiosArgs.options, url: (configuration?.basePath || basePath) + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
        /**
         * 
         * @summary Updates SpellBucket
         * @param {number} id Id
         * @param {ModelsSpellBucket} spellBucket SpellBucket
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async updateSpellBucket(id: number, spellBucket: ModelsSpellBucket, options?: any): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<Array<ModelsSpellBucket>>> {
            const localVarAxiosArgs = await SpellBucketApiAxiosParamCreator(configuration).updateSpellBucket(id, spellBucket, options);
            return (axios: AxiosInstance = globalAxios, basePath: string = BASE_PATH) => {
                const axiosRequestArgs = {...localVarAxiosArgs.options, url: (configuration?.basePath || basePath) + localVarAxiosArgs.url};
                return axios.request(axiosRequestArgs);
            };
        },
    }
};

/**
 * SpellBucketApi - factory interface
 * @export
 */
export const SpellBucketApiFactory = function (configuration?: Configuration, basePath?: string, axios?: AxiosInstance) {
    return {
        /**
         * 
         * @summary Creates SpellBucket
         * @param {ModelsSpellBucket} spellBucket SpellBucket
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        createSpellBucket(spellBucket: ModelsSpellBucket, options?: any): AxiosPromise<Array<ModelsSpellBucket>> {
            return SpellBucketApiFp(configuration).createSpellBucket(spellBucket, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Deletes SpellBucket
         * @param {number} id spellid
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        deleteSpellBucket(id: number, options?: any): AxiosPromise<string> {
            return SpellBucketApiFp(configuration).deleteSpellBucket(id, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Gets SpellBucket
         * @param {number} id Id
         * @param {string} [includes] Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names 
         * @param {string} [select] Column names [.] separated to fetch specific fields in response
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getSpellBucket(id: number, includes?: string, select?: string, options?: any): AxiosPromise<Array<ModelsSpellBucket>> {
            return SpellBucketApiFp(configuration).getSpellBucket(id, includes, select, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Gets SpellBuckets in bulk
         * @param {CrudcontrollersBulkFetchByIdsGetRequest} body body
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getSpellBucketsBulk(body: CrudcontrollersBulkFetchByIdsGetRequest, options?: any): AxiosPromise<Array<ModelsSpellBucket>> {
            return SpellBucketApiFp(configuration).getSpellBucketsBulk(body, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Lists SpellBuckets
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
        listSpellBuckets(includes?: string, where?: string, whereOr?: string, groupBy?: string, limit?: string, orderBy?: string, orderDirection?: string, select?: string, options?: any): AxiosPromise<Array<ModelsSpellBucket>> {
            return SpellBucketApiFp(configuration).listSpellBuckets(includes, where, whereOr, groupBy, limit, orderBy, orderDirection, select, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Updates SpellBucket
         * @param {number} id Id
         * @param {ModelsSpellBucket} spellBucket SpellBucket
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        updateSpellBucket(id: number, spellBucket: ModelsSpellBucket, options?: any): AxiosPromise<Array<ModelsSpellBucket>> {
            return SpellBucketApiFp(configuration).updateSpellBucket(id, spellBucket, options).then((request) => request(axios, basePath));
        },
    };
};

/**
 * Request parameters for createSpellBucket operation in SpellBucketApi.
 * @export
 * @interface SpellBucketApiCreateSpellBucketRequest
 */
export interface SpellBucketApiCreateSpellBucketRequest {
    /**
     * SpellBucket
     * @type {ModelsSpellBucket}
     * @memberof SpellBucketApiCreateSpellBucket
     */
    readonly spellBucket: ModelsSpellBucket
}

/**
 * Request parameters for deleteSpellBucket operation in SpellBucketApi.
 * @export
 * @interface SpellBucketApiDeleteSpellBucketRequest
 */
export interface SpellBucketApiDeleteSpellBucketRequest {
    /**
     * spellid
     * @type {number}
     * @memberof SpellBucketApiDeleteSpellBucket
     */
    readonly id: number
}

/**
 * Request parameters for getSpellBucket operation in SpellBucketApi.
 * @export
 * @interface SpellBucketApiGetSpellBucketRequest
 */
export interface SpellBucketApiGetSpellBucketRequest {
    /**
     * Id
     * @type {number}
     * @memberof SpellBucketApiGetSpellBucket
     */
    readonly id: number

    /**
     * Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names 
     * @type {string}
     * @memberof SpellBucketApiGetSpellBucket
     */
    readonly includes?: string

    /**
     * Column names [.] separated to fetch specific fields in response
     * @type {string}
     * @memberof SpellBucketApiGetSpellBucket
     */
    readonly select?: string
}

/**
 * Request parameters for getSpellBucketsBulk operation in SpellBucketApi.
 * @export
 * @interface SpellBucketApiGetSpellBucketsBulkRequest
 */
export interface SpellBucketApiGetSpellBucketsBulkRequest {
    /**
     * body
     * @type {CrudcontrollersBulkFetchByIdsGetRequest}
     * @memberof SpellBucketApiGetSpellBucketsBulk
     */
    readonly body: CrudcontrollersBulkFetchByIdsGetRequest
}

/**
 * Request parameters for listSpellBuckets operation in SpellBucketApi.
 * @export
 * @interface SpellBucketApiListSpellBucketsRequest
 */
export interface SpellBucketApiListSpellBucketsRequest {
    /**
     * Relationships [all] for all [number] for depth of relationships to load or [.] separated relationship names 
     * @type {string}
     * @memberof SpellBucketApiListSpellBuckets
     */
    readonly includes?: string

    /**
     * Filter on specific fields. Multiple conditions [.] separated Example: col_like_value.col2__val2
     * @type {string}
     * @memberof SpellBucketApiListSpellBuckets
     */
    readonly where?: string

    /**
     * Filter on specific fields (Chained ors). Multiple conditions [.] separated Example: col_like_value.col2__val2
     * @type {string}
     * @memberof SpellBucketApiListSpellBuckets
     */
    readonly whereOr?: string

    /**
     * Group by field. Multiple conditions [.] separated Example: field1.field2
     * @type {string}
     * @memberof SpellBucketApiListSpellBuckets
     */
    readonly groupBy?: string

    /**
     * Rows to limit in response (Default: 10,000)
     * @type {string}
     * @memberof SpellBucketApiListSpellBuckets
     */
    readonly limit?: string

    /**
     * Order by [field]
     * @type {string}
     * @memberof SpellBucketApiListSpellBuckets
     */
    readonly orderBy?: string

    /**
     * Order by field direction
     * @type {string}
     * @memberof SpellBucketApiListSpellBuckets
     */
    readonly orderDirection?: string

    /**
     * Column names [.] separated to fetch specific fields in response
     * @type {string}
     * @memberof SpellBucketApiListSpellBuckets
     */
    readonly select?: string
}

/**
 * Request parameters for updateSpellBucket operation in SpellBucketApi.
 * @export
 * @interface SpellBucketApiUpdateSpellBucketRequest
 */
export interface SpellBucketApiUpdateSpellBucketRequest {
    /**
     * Id
     * @type {number}
     * @memberof SpellBucketApiUpdateSpellBucket
     */
    readonly id: number

    /**
     * SpellBucket
     * @type {ModelsSpellBucket}
     * @memberof SpellBucketApiUpdateSpellBucket
     */
    readonly spellBucket: ModelsSpellBucket
}

/**
 * SpellBucketApi - object-oriented interface
 * @export
 * @class SpellBucketApi
 * @extends {BaseAPI}
 */
export class SpellBucketApi extends BaseAPI {
    /**
     * 
     * @summary Creates SpellBucket
     * @param {SpellBucketApiCreateSpellBucketRequest} requestParameters Request parameters.
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof SpellBucketApi
     */
    public createSpellBucket(requestParameters: SpellBucketApiCreateSpellBucketRequest, options?: any) {
        return SpellBucketApiFp(this.configuration).createSpellBucket(requestParameters.spellBucket, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Deletes SpellBucket
     * @param {SpellBucketApiDeleteSpellBucketRequest} requestParameters Request parameters.
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof SpellBucketApi
     */
    public deleteSpellBucket(requestParameters: SpellBucketApiDeleteSpellBucketRequest, options?: any) {
        return SpellBucketApiFp(this.configuration).deleteSpellBucket(requestParameters.id, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Gets SpellBucket
     * @param {SpellBucketApiGetSpellBucketRequest} requestParameters Request parameters.
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof SpellBucketApi
     */
    public getSpellBucket(requestParameters: SpellBucketApiGetSpellBucketRequest, options?: any) {
        return SpellBucketApiFp(this.configuration).getSpellBucket(requestParameters.id, requestParameters.includes, requestParameters.select, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Gets SpellBuckets in bulk
     * @param {SpellBucketApiGetSpellBucketsBulkRequest} requestParameters Request parameters.
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof SpellBucketApi
     */
    public getSpellBucketsBulk(requestParameters: SpellBucketApiGetSpellBucketsBulkRequest, options?: any) {
        return SpellBucketApiFp(this.configuration).getSpellBucketsBulk(requestParameters.body, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Lists SpellBuckets
     * @param {SpellBucketApiListSpellBucketsRequest} requestParameters Request parameters.
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof SpellBucketApi
     */
    public listSpellBuckets(requestParameters: SpellBucketApiListSpellBucketsRequest = {}, options?: any) {
        return SpellBucketApiFp(this.configuration).listSpellBuckets(requestParameters.includes, requestParameters.where, requestParameters.whereOr, requestParameters.groupBy, requestParameters.limit, requestParameters.orderBy, requestParameters.orderDirection, requestParameters.select, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Updates SpellBucket
     * @param {SpellBucketApiUpdateSpellBucketRequest} requestParameters Request parameters.
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof SpellBucketApi
     */
    public updateSpellBucket(requestParameters: SpellBucketApiUpdateSpellBucketRequest, options?: any) {
        return SpellBucketApiFp(this.configuration).updateSpellBucket(requestParameters.id, requestParameters.spellBucket, options).then((request) => request(this.axios, this.basePath));
    }
}