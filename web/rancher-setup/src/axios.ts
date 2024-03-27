// axios.ts

import axios, { AxiosInstance, AxiosResponse } from 'axios';

export class MyHttp {
  private axiosInstance: AxiosInstance;

  constructor(baseURL: string) {
    this.axiosInstance = axios.create({
      baseURL,
      timeout: 20 * 60 * 1000, // 设置超时时间
    });
  }

  // 封装 get 请求
  async get<T = any>(url: string): Promise<AxiosResponse<T>> {
    return this.axiosInstance.get<T>(url);
  }

  // 封装 post 请求
  async post<T = any>(url: string, data: any): Promise<AxiosResponse<T>> {
    return this.axiosInstance.post<T>(url, data);
  }

  // 封装其他 HTTP 方法...
}