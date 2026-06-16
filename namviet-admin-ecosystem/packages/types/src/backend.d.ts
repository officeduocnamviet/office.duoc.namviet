/* eslint-disable */
/* tslint:disable */
// @ts-nocheck
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface BatchesBatch {
  batch_code?: string;
  created_at?: string;
  deleted_at?: string;
  expiry_date?: string;
  id?: number;
  inbound_price?: number;
  manufacturing_date?: string;
  product_id?: number;
  updated_at?: string;
}

export interface BatchesCreateBatchRequest {
  batch_code: string;
  expiry_date: string;
  inbound_price?: number;
  manufacturing_date?: string;
  product_id: number;
}

export interface BatchesUpdateBatchRequest {
  batch_code?: string;
  expiry_date?: string;
  inbound_price?: number;
  manufacturing_date?: string;
}

export interface CategoriesCategory {
  created_at?: string;
  deleted_at?: string;
  id?: number;
  name?: string;
  parent_id?: number;
  slug?: string;
  status?: string;
  updated_at?: string;
}

export interface CategoriesCreateCategoryRequest {
  name: string;
  parent_id?: number;
  slug: string;
  status?: string;
}

export interface CategoriesUpdateCategoryRequest {
  name?: string;
  parent_id?: number;
  slug?: string;
  status?: string;
}

export interface InventoryCreateTransactionRequest {
  batch_id?: number;
  product_id: number;
  quantity: number;
  reference_id?: string;
  reference_type?: string;
  /** IN, OUT */
  type: string;
  warehouse_id: number;
}

export interface InventoryInventoryBatch {
  batch_id?: number;
  created_at?: string;
  id?: number;
  product_id?: number;
  quantity?: number;
  updated_at?: string;
  warehouse_id?: number;
}

export interface InventoryInventoryTransaction {
  batch_id?: number;
  created_at?: string;
  id?: number;
  product_id?: number;
  quantity?: number;
  reference_id?: string;
  reference_type?: string;
  /** e.g. IN, OUT */
  type?: string;
  warehouse_id?: number;
}

export interface ManufacturersCreateManufacturerRequest {
  country?: string;
  name: string;
  status?: string;
}

export interface ManufacturersManufacturer {
  country?: string;
  created_at?: string;
  deleted_at?: string;
  id?: number;
  name?: string;
  status?: string;
  updated_at?: string;
}

export interface ManufacturersUpdateManufacturerRequest {
  country?: string;
  name?: string;
  status?: string;
}

export interface ProductUnitsCreateProductUnitRequest {
  conversion_factor: number;
  is_base_unit?: boolean;
  price_cost?: number;
  price_sell?: number;
  unit_name: string;
}

export interface ProductUnitsProductUnit {
  conversion_factor?: number;
  created_at?: string;
  deleted_at?: string;
  id?: number;
  is_base_unit?: boolean;
  price_cost?: number;
  price_sell?: number;
  product_id?: number;
  unit_name?: string;
  updated_at?: string;
}

export interface ProductUnitsUpdateProductUnitRequest {
  conversion_factor?: number;
  is_base_unit?: boolean;
  price_cost?: number;
  price_sell?: number;
  unit_name?: string;
}

export interface ProductsCreateProductRequest {
  actual_cost?: number;
  barcode?: string;
  category_id?: number;
  category_name?: string;
  conversion_factor?: number;
  description?: string;
  manufacturer_id?: number;
  manufacturer_name?: string;
  name: string;
  price_cost?: number;
  price_sell?: number;
  retail_unit?: string;
  sku?: string;
  wholesale_unit?: string;
}

export interface ProductsProduct {
  active_ingredient?: string;
  actual_cost?: number;
  barcode?: string;
  carton_dimensions?: string;
  carton_weight?: number;
  category_id?: number;
  category_name?: string;
  conversion_factor?: number;
  created_at?: string;
  deleted_at?: string;
  description?: string;
  distributor_id?: number;
  id?: number;
  image_url?: string;
  invoice_price?: number;
  is_direct_sale?: boolean;
  items_per_carton?: number;
  manufacturer_id?: number;
  manufacturer_name?: string;
  name?: string;
  price?: number;
  price_cost?: number;
  price_sell?: number;
  purchasing_policy?: string;
  registration_number?: string;
  retail_margin_type?: string;
  retail_margin_value?: number;
  retail_unit?: string;
  sku?: string;
  status?: string;
  unit_type?: string;
  updated_at?: string;
  wholesale_margin_type?: string;
  wholesale_margin_value?: number;
  wholesale_unit?: string;
}

export interface ProductsUpdateProductRequest {
  actual_cost?: number;
  barcode?: string;
  category_id?: number;
  category_name?: string;
  description?: string;
  manufacturer_id?: number;
  manufacturer_name?: string;
  name?: string;
  sku?: string;
  status?: string;
}

export interface RolesCreateRoleRequest {
  description?: string;
  name: string;
  permissions?: string[];
}

export interface RolesRole {
  created_at?: string;
  description?: string;
  id?: string;
  name?: string;
  permissions?: string[];
}

export interface RolesUpdateRoleRequest {
  description?: string;
  name?: string;
  permissions?: string[];
}

export interface UsersCreateUserRequest {
  company_id: string;
  email: string;
  full_name: string;
  /** @minLength 6 */
  password: string;
  phone?: string;
  role_id: string;
  warehouse_id?: number;
}

export interface UsersLoginRequest {
  /** @example "admin@namviet.com" */
  email: string;
  /** @example "namviet123" */
  password: string;
}

export interface UsersLoginResponse {
  token?: string;
  user?: UsersUser;
}

export interface UsersRegisterFCMTokenRequest {
  device_info?: string;
  token: string;
}

export interface UsersUpdateUserRequest {
  full_name?: string;
  phone?: string;
  role_id?: string;
  status?: string;
  warehouse_id?: number;
}

export interface UsersUser {
  company_id?: string;
  created_at?: string;
  deleted_at?: string;
  email?: string;
  full_name?: string;
  id?: string;
  phone?: string;
  role_id?: string;
  status?: string;
  updated_at?: string;
  warehouse_id?: number;
}

export interface WarehousesCreateWarehouseRequest {
  address?: string;
  manager_id?: string;
  name: string;
  status?: string;
  type?: string;
}

export interface WarehousesUpdateWarehouseRequest {
  address?: string;
  manager_id?: string;
  name?: string;
  status?: string;
  type?: string;
}

export interface WarehousesWarehouse {
  address?: string;
  created_at?: string;
  deleted_at?: string;
  id?: number;
  manager_id?: string;
  name?: string;
  status?: string;
  type?: string;
  updated_at?: string;
}
