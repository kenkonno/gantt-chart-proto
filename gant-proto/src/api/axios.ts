import {DefaultApiFactory} from "@/api/api";
import axios, {CreateAxiosDefaults} from "axios";
import {Configuration, ConfigurationParameters} from "@/api/configuration";

const param: ConfigurationParameters = {}
const configuration = new Configuration(param);
const basePath = ""
const axiosConfig: CreateAxiosDefaults = {
    baseURL: "",
    headers: {
        // 'Access-Control-Allow-Origin': '*',
        // 'Access-Control-Allow-Headers': '*',
        // 'Access-Control-Allow-Credentials': 'true',
        // 'Content-Type': 'text/plain'
        // 何だったんだろうかこれは？・・・・
    },
    withCredentials: true
}
const axiosInstance = axios.create(axiosConfig)
export const Api = DefaultApiFactory(configuration, basePath, axiosInstance)