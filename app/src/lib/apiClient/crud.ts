import type {
  ApiResponse,
  CreateResponse,
  GetResponse,
  UpdateResponse,
  DeleteResponse,
} from "./types"

const API_BASE = "/api/web/crud"

function getToken(): string {
  const token = localStorage.getItem("sessionToken")
  if (!token)
    throw new Error("セッショントークンがありません。ログインしてください。")
  return token
}

async function post<T>(
  endpoint: string,
  body: object,
): Promise<ApiResponse<T>> {
  const res = await fetch(`${API_BASE}${endpoint}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${getToken()}`,
    },
    body: JSON.stringify(body),
  })
  return res.json() as Promise<ApiResponse<T>>
}

export function createData(
  key: string,
  value: string,
): Promise<ApiResponse<CreateResponse>> {
  return post<CreateResponse>("/create", { key, value })
}

export function getData(key: string): Promise<ApiResponse<GetResponse>> {
  return post<GetResponse>("/get", { key })
}

export function updateData(
  key: string,
  value: string,
): Promise<ApiResponse<UpdateResponse>> {
  return post<UpdateResponse>("/update", { key, value })
}

export function deleteData(key: string): Promise<ApiResponse<DeleteResponse>> {
  return post<DeleteResponse>("/delete", { key })
}
