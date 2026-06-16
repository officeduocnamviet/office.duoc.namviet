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

export interface AccountingJournalsAccountingJournal {
  account_credit?: string;
  account_debit?: string;
  amount?: number;
  created_at?: string;
  description?: string;
  doc_type?: string;
  entry_date?: string;
  id?: string;
  posted_by?: string;
  source_ref_id?: string;
}

export interface AccountingJournalsCreateAccountingJournalRequest {
  account_credit: string;
  account_debit: string;
  amount: number;
  description?: string;
  doc_type: string;
  entry_date: string;
  posted_by?: string;
  source_ref_id?: string;
}

export interface AccountingJournalsUpdateAccountingJournalRequest {
  account_credit?: string;
  account_debit?: string;
  amount?: number;
  description?: string;
  doc_type?: string;
  entry_date?: string;
  posted_by?: string;
  source_ref_id?: string;
}

export interface AgentWorkflowsAgentWorkflow {
  created_at?: string;
  description?: string;
  id?: string;
  is_active?: boolean;
  name?: string;
  steps?: string[];
  trigger_type?: string;
  updated_at?: string;
}

export interface AgentWorkflowsCreateAgentWorkflowRequest {
  description?: string;
  is_active?: boolean;
  name: string;
  steps: string[];
  trigger_type: string;
}

export interface AgentWorkflowsUpdateAgentWorkflowRequest {
  description?: string;
  is_active?: boolean;
  name?: string;
  steps?: string[];
  trigger_type?: string;
}

export interface AiAgentMemoriesAIAgentMemory {
  created_at?: string;
  expires_at?: string;
  id?: string;
  key?: string;
  session_id?: string;
  user_id?: string;
  value?: string[];
}

export interface AiAgentMemoriesCreateAIAgentMemoryRequest {
  expires_at?: string;
  key: string;
  session_id?: string;
  user_id?: string;
  value: string[];
}

export interface AiAgentMemoriesUpdateAIAgentMemoryRequest {
  expires_at?: string;
  key?: string;
  session_id?: string;
  user_id?: string;
  value?: string[];
}

export interface AppointmentsAppointment {
  appointment_time?: string;
  check_in_time?: string;
  created_at?: string;
  created_by?: string;
  customer_id?: number;
  deleted_at?: string;
  doctor_id?: string;
  id?: string;
  note?: string;
  room_id?: number;
  service_type?: string;
  status?: string;
  symptoms?: string[];
  updated_at?: string;
}

export interface AppointmentsCreateAppointmentRequest {
  appointment_time: string;
  customer_id: number;
  doctor_id?: string;
  note?: string;
  room_id?: number;
  service_type: string;
  symptoms?: string[];
}

export interface AppointmentsUpdateAppointmentRequest {
  appointment_time?: string;
  check_in_time?: string;
  doctor_id?: string;
  note?: string;
  room_id?: number;
  service_type?: string;
  status?: string;
  symptoms?: string[];
}

export interface ApprovalsApprovalRequest {
  created_at?: string;
  current_step?: number;
  id?: string;
  payload?: string[];
  ref_id?: string;
  request_type?: string;
  requester_id?: string;
  status?: string;
  updated_at?: string;
}

export interface ApprovalsApprovalStep {
  action_at?: string;
  approver_id?: string;
  approver_role?: string;
  comments?: string;
  created_at?: string;
  id?: string;
  request_id?: string;
  status?: string;
  step_order?: number;
}

export interface ApprovalsCreateApprovalRequestDto {
  payload?: string[];
  ref_id?: string;
  request_type: string;
  requester_id?: string;
}

export interface ApprovalsCreateApprovalStepDto {
  approver_id?: string;
  approver_role?: string;
  request_id?: string;
  step_order: number;
}

export interface ApprovalsUpdateApprovalRequestDto {
  current_step?: number;
  payload?: string[];
  status?: string;
}

export interface ApprovalsUpdateApprovalStepDto {
  comments?: string;
  status?: string;
}

export interface AttendanceLogsAttendanceLog {
  branch_id?: number;
  check_in_ip?: string;
  check_in_lat?: number;
  check_in_lng?: number;
  check_in_time?: string;
  check_out_ip?: string;
  check_out_lat?: number;
  check_out_lng?: number;
  check_out_time?: string;
  id?: string;
  status?: string;
  user_id?: string;
}

export interface AttendanceLogsCreateAttendanceLogRequest {
  branch_id?: number;
  check_in_ip?: string;
  check_in_lat?: number;
  check_in_lng?: number;
  check_in_time: string;
  user_id: string;
}

export interface AttendanceLogsUpdateAttendanceLogRequest {
  check_out_ip?: string;
  check_out_lat?: number;
  check_out_lng?: number;
  check_out_time?: string;
  status?: string;
}

export interface AuditLogsSystemAuditLog {
  action?: string;
  created_at?: string;
  id?: string;
  ip_address?: string;
  new_data?: string[];
  old_data?: string[];
  record_id?: string;
  table_name?: string;
  user_agent?: string;
  user_id?: string;
}

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

export interface ChartOfAccountsChartOfAccount {
  account_code?: string;
  allow_posting?: boolean;
  balance_type?: string;
  created_at?: string;
  deleted_at?: string;
  id?: string;
  name?: string;
  parent_id?: string;
  status?: string;
  type?: string;
  updated_at?: string;
}

export interface ChartOfAccountsCreateChartOfAccountRequest {
  account_code: string;
  allow_posting?: boolean;
  balance_type: string;
  name: string;
  parent_id?: string;
  type: string;
}

export interface ChartOfAccountsUpdateChartOfAccountRequest {
  account_code?: string;
  allow_posting?: boolean;
  balance_type?: string;
  name?: string;
  parent_id?: string;
  status?: string;
  type?: string;
}

export interface ChatsChatMessage {
  content?: string;
  created_at?: string;
  id?: string;
  message_type?: string;
  metadata?: string[];
  sender_id?: string;
  sender_type?: string;
  session_id?: string;
}

export interface ChatsChatSession {
  agent_id?: string;
  created_at?: string;
  customer_id?: string;
  ended_at?: string;
  id?: string;
  status?: string;
}

export interface ChatsCreateChatMessageRequest {
  content: string;
  message_type?: string;
  metadata?: string[];
  sender_id?: string;
  sender_type: string;
  session_id?: string;
}

export interface ChatsCreateChatSessionRequest {
  agent_id?: string;
  customer_id?: string;
}

export interface ChatsUpdateChatSessionRequest {
  status?: string;
}

export interface ClinicalQueuesClinicalQueue {
  appointment_id?: string;
  checked_in_at?: string;
  customer_id?: number;
  deleted_at?: string;
  doctor_id?: string;
  id?: string;
  priority_level?: string;
  queue_number?: number;
  status?: string;
  updated_at?: string;
}

export interface ClinicalQueuesCreateClinicalQueueRequest {
  appointment_id?: string;
  customer_id: number;
  doctor_id?: string;
  priority_level?: string;
  queue_number: number;
}

export interface ClinicalQueuesUpdateClinicalQueueRequest {
  doctor_id?: string;
  priority_level?: string;
  status?: string;
}

export interface CompaniesBranch {
  address?: string;
  code?: string;
  company_id?: string;
  created_at?: string;
  deleted_at?: string;
  id?: string;
  latitude?: number;
  longitude?: number;
  manager_id?: string;
  name?: string;
  phone?: string;
  status?: string;
  updated_at?: string;
}

export interface CompaniesCompany {
  address?: string;
  business_image_license_url?: string[];
  created_at?: string;
  deleted_at?: string;
  email?: string;
  id?: string;
  logo_url?: string;
  name?: string;
  phone?: string;
  representative_name?: string;
  short_name?: string;
  status?: string;
  tax_code?: string;
  updated_at?: string;
}

export interface CompaniesCreateBranchRequest {
  address: string;
  code: string;
  company_id: string;
  latitude?: number;
  longitude?: number;
  manager_id?: string;
  name: string;
  phone?: string;
}

export interface CompaniesCreateCompanyRequest {
  address: string;
  business_image_license_url?: string[];
  email?: string;
  logo_url?: string;
  name: string;
  phone: string;
  representative_name?: string;
  short_name?: string;
  tax_code: string;
}

export interface CompaniesUpdateBranchRequest {
  address?: string;
  code?: string;
  company_id?: string;
  latitude?: number;
  longitude?: number;
  manager_id?: string;
  name?: string;
  phone?: string;
  status?: string;
}

export interface CompaniesUpdateCompanyRequest {
  address?: string;
  business_image_license_url?: string[];
  email?: string;
  logo_url?: string;
  name?: string;
  phone?: string;
  representative_name?: string;
  short_name?: string;
  status?: string;
  tax_code?: string;
}

export interface CustomerRecordsCreateCustomerVoucherRequest {
  customer_id: string;
  promotion_id: number;
  voucher_code: string;
}

export interface CustomerRecordsCreateVaccinationRecordRequest {
  administered_by?: string;
  customer_id: string;
  dose_number?: number;
  next_due_date?: string;
  notes?: string;
  vaccination_date: string;
  vaccine_name: string;
}

export interface CustomerRecordsCustomerVaccinationRecord {
  administered_by?: string;
  created_at?: string;
  customer_id?: string;
  dose_number?: number;
  id?: string;
  next_due_date?: string;
  notes?: string;
  vaccination_date?: string;
  vaccine_name?: string;
}

export interface CustomerRecordsCustomerVoucher {
  created_at?: string;
  customer_id?: string;
  id?: string;
  is_used?: boolean;
  promotion_id?: number;
  used_at?: string;
  voucher_code?: string;
}

export interface CustomerRecordsUpdateCustomerVoucherRequest {
  is_used?: boolean;
  used_at?: string;
}

export interface CustomerRecordsUpdateVaccinationRecordRequest {
  administered_by?: string;
  dose_number?: number;
  next_due_date?: string;
  notes?: string;
  vaccination_date?: string;
  vaccine_name?: string;
}

export interface CustomersCreateCustomerRequest {
  address?: string;
  b2b_metadata?: string[];
  cccd?: string;
  customer_code?: string;
  /** B2B or B2C */
  customer_type?: string;
  dob?: string;
  email?: string;
  gender?: string;
  name: string;
  phone?: string;
}

export interface CustomersCustomer {
  address?: string;
  b2b_metadata?: string[];
  cccd?: string;
  created_at?: string;
  current_debt?: number;
  customer_code?: string;
  customer_type?: string;
  deleted_at?: string;
  dob?: string;
  email?: string;
  gender?: string;
  id?: number;
  loyalty_points?: number;
  name?: string;
  phone?: string;
  status?: string;
  updated_at?: string;
  updated_by?: string;
}

export interface CustomersUpdateCustomerRequest {
  address?: string;
  b2b_metadata?: string[];
  cccd?: string;
  customer_code?: string;
  customer_type?: string;
  dob?: string;
  email?: string;
  gender?: string;
  name?: string;
  phone?: string;
  status?: string;
}

export interface EmployeesCreateEmployeeRequest {
  bank_account?: string[];
  base_salary?: number;
  code?: string;
  department_id?: string;
  hire_date?: string;
  insurance_no?: string;
  position_id?: string;
  salary_type?: string;
  tax_code?: string;
  type?: string;
  user_id?: string;
}

export interface EmployeesEmployee {
  bank_account?: string[];
  base_salary?: number;
  code?: string;
  created_at?: string;
  deleted_at?: string;
  department_id?: string;
  hire_date?: string;
  id?: string;
  insurance_no?: string;
  position_id?: string;
  salary_type?: string;
  status?: string;
  tax_code?: string;
  termination_date?: string;
  type?: string;
  updated_at?: string;
  user_id?: string;
}

export interface EmployeesUpdateEmployeeRequest {
  bank_account?: string[];
  base_salary?: number;
  code?: string;
  department_id?: string;
  hire_date?: string;
  insurance_no?: string;
  position_id?: string;
  salary_type?: string;
  status?: string;
  tax_code?: string;
  termination_date?: string;
  type?: string;
}

export interface EmploymentContractsCreateEmploymentContractRequest {
  base_salary?: number;
  commission_rate_percent?: number;
  contract_code: string;
  insurance_deduction_amount?: number;
  kpi_conversion_rate?: number;
  standard_working_days?: number;
  tax_deduction_amount?: number;
  user_id: string;
  valid_from: string;
  valid_to?: string;
}

export interface EmploymentContractsEmploymentContract {
  base_salary?: number;
  commission_rate_percent?: number;
  contract_code?: string;
  created_at?: string;
  id?: string;
  insurance_deduction_amount?: number;
  kpi_conversion_rate?: number;
  standard_working_days?: number;
  status?: string;
  tax_deduction_amount?: number;
  updated_at?: string;
  user_id?: string;
  valid_from?: string;
  valid_to?: string;
}

export interface EmploymentContractsUpdateEmploymentContractRequest {
  base_salary?: number;
  commission_rate_percent?: number;
  contract_code?: string;
  insurance_deduction_amount?: number;
  kpi_conversion_rate?: number;
  standard_working_days?: number;
  status?: string;
  tax_deduction_amount?: number;
  valid_from?: string;
  valid_to?: string;
}

export interface FinanceTransactionsCreateFinanceTransactionRequest {
  amount: number;
  bank_reference_id?: string;
  book_type?: string;
  business_type?: string;
  cash_tally?: string[];
  category_id?: number;
  code: string;
  created_by?: string;
  description?: string;
  evidence_url?: string;
  flow: string;
  fund_account_id: number;
  partner_id?: string;
  partner_name_cache?: string;
  partner_type?: string;
  ref_advance_id?: number;
  ref_id?: string;
  ref_type?: string;
  target_bank_info?: string[];
  transaction_date?: string;
}

export interface FinanceTransactionsFinanceTransaction {
  amount?: number;
  bank_reference_id?: string;
  book_type?: string;
  business_type?: string;
  cash_tally?: string[];
  category_id?: number;
  code?: string;
  created_at?: string;
  created_by?: string;
  deleted_at?: string;
  description?: string;
  evidence_url?: string;
  flow?: string;
  fund_account_id?: number;
  id?: number;
  is_posted?: boolean;
  partner_id?: string;
  partner_name_cache?: string;
  partner_type?: string;
  ref_advance_id?: number;
  ref_id?: string;
  ref_type?: string;
  status?: string;
  target_bank_info?: string[];
  transaction_date?: string;
  updated_at?: string;
}

export interface FinanceTransactionsUpdateFinanceTransactionRequest {
  amount?: number;
  bank_reference_id?: string;
  book_type?: string;
  business_type?: string;
  cash_tally?: string[];
  category_id?: number;
  description?: string;
  evidence_url?: string;
  flow?: string;
  fund_account_id?: number;
  is_posted?: boolean;
  partner_id?: string;
  partner_name_cache?: string;
  partner_type?: string;
  ref_advance_id?: number;
  ref_id?: string;
  ref_type?: string;
  status?: string;
  target_bank_info?: string[];
  transaction_date?: string;
}

export interface FundAccountsCreateFundAccountRequest {
  account_id?: string;
  account_number?: string;
  balance?: number;
  bank_id?: number;
  bank_info?: string[];
  currency?: string;
  description?: string;
  initial_balance?: number;
  location?: string;
  name: string;
  type: string;
}

export interface FundAccountsFundAccount {
  account_id?: string;
  account_number?: string;
  balance?: number;
  bank_id?: number;
  bank_info?: string[];
  created_at?: string;
  currency?: string;
  deleted_at?: string;
  description?: string;
  id?: number;
  initial_balance?: number;
  location?: string;
  name?: string;
  status?: string;
  type?: string;
  updated_at?: string;
}

export interface FundAccountsUpdateFundAccountRequest {
  account_id?: string;
  account_number?: string;
  balance?: number;
  bank_id?: number;
  bank_info?: string[];
  currency?: string;
  description?: string;
  location?: string;
  name?: string;
  status?: string;
  type?: string;
}

export interface GormDeletedAt {
  time?: string;
  /** Valid is true if Time is not NULL */
  valid?: boolean;
}

export interface IntegrationsCreateConnectionRequest {
  api_key?: string;
  partner_name: string;
  secret_key?: string;
  webhook_url?: string;
}

export interface IntegrationsCreateWebhookLogRequest {
  event_type: string;
  partner_id?: string;
  payload?: string[];
  response_body?: string;
  response_status?: number;
}

export interface IntegrationsThirdPartyConnection {
  api_key?: string;
  created_at?: string;
  id?: string;
  partner_name?: string;
  secret_key?: string;
  status?: string;
  updated_at?: string;
  webhook_url?: string;
}

export interface IntegrationsUpdateConnectionRequest {
  api_key?: string;
  partner_name?: string;
  secret_key?: string;
  status?: string;
  webhook_url?: string;
}

export interface IntegrationsWebhookLog {
  created_at?: string;
  event_type?: string;
  id?: string;
  partner_id?: string;
  payload?: string[];
  response_body?: string;
  response_status?: number;
}

export interface InternalCommunicationsCreateInternalChannelRequest {
  name: string;
  type?: string;
}

export interface InternalCommunicationsCreateInternalMessageRequest {
  channel_id: number;
  content: string;
  sender_id: string;
}

export interface InternalCommunicationsInternalChannel {
  created_at?: string;
  id?: number;
  name?: string;
  type?: string;
}

export interface InternalCommunicationsInternalMessage {
  channel_id?: number;
  content?: string;
  created_at?: string;
  id?: number;
  sender_id?: string;
}

export interface InternalCommunicationsUpdateInternalChannelRequest {
  name?: string;
  type?: string;
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

export interface KnowledgeVectorsCreateMedicalKnowledgeVectorRequest {
  content: string;
  embedding?: string;
  metadata?: string[];
  title: string;
}

export interface KnowledgeVectorsCreateProductVectorRequest {
  content: string;
  embedding?: string;
  metadata?: string[];
  product_id?: string;
}

export interface KnowledgeVectorsMedicalKnowledgeVector {
  content?: string;
  created_at?: string;
  /** simplified */
  embedding?: string;
  id?: string;
  metadata?: string[];
  title?: string;
}

export interface KnowledgeVectorsProductVector {
  content?: string;
  created_at?: string;
  /** simplified */
  embedding?: string;
  id?: string;
  metadata?: string[];
  product_id?: string;
}

export interface KnowledgeVectorsUpdateMedicalKnowledgeVectorRequest {
  content?: string;
  embedding?: string;
  metadata?: string[];
  title?: string;
}

export interface KnowledgeVectorsUpdateProductVectorRequest {
  content?: string;
  embedding?: string;
  metadata?: string[];
  product_id?: string;
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

export interface MarketingCampaignsCreateMarketingCampaignRequest {
  budget?: number;
  description?: string;
  end_date?: string;
  name: string;
  start_date: string;
  target_segment?: string[];
}

export interface MarketingCampaignsMarketingCampaign {
  budget?: number;
  created_at?: string;
  description?: string;
  end_date?: string;
  id?: string;
  name?: string;
  start_date?: string;
  status?: string;
  target_segment?: string[];
  updated_at?: string;
}

export interface MarketingCampaignsUpdateMarketingCampaignRequest {
  budget?: number;
  description?: string;
  end_date?: string;
  name?: string;
  start_date?: string;
  status?: string;
  target_segment?: string[];
}

export interface MedicalVisitsCreateMedicalVisitRequest {
  appointment_id?: string;
  customer_id: number;
  doctor_id?: string;
  symptoms?: string;
}

export interface MedicalVisitsMedicalVisit {
  appointment_id?: string;
  bp_diastolic?: number;
  bp_systolic?: number;
  created_at?: string;
  created_by?: string;
  customer_id?: number;
  deleted_at?: string;
  diagnosis?: string;
  doctor_id?: string;
  doctor_notes?: string;
  examination_summary?: string;
  height?: number;
  icd_code?: string;
  id?: string;
  pulse?: number;
  red_flags?: string[];
  sp02?: number;
  status?: string;
  symptoms?: string;
  temperature?: number;
  updated_at?: string;
  updated_by?: string;
  weight?: number;
}

export interface MedicalVisitsUpdateMedicalVisitRequest {
  bp_diastolic?: number;
  bp_systolic?: number;
  diagnosis?: string;
  doctor_notes?: string;
  examination_summary?: string;
  height?: number;
  icd_code?: string;
  pulse?: number;
  red_flags?: string[];
  sp02?: number;
  status?: string;
  symptoms?: string;
  temperature?: number;
  weight?: number;
}

export interface OrdersCreateOrderItemRequest {
  batch_no?: string;
  conversion_factor?: number;
  discount?: number;
  expiry_date?: string;
  is_gift?: boolean;
  note?: string;
  product_id: number;
  quantity: number;
  unit_price: number;
  uom: string;
}

export interface OrdersCreateOrderRequest {
  code: string;
  creator_id?: string;
  customer_id?: number;
  items: OrdersCreateOrderItemRequest[];
  note?: string;
  /** B2B or B2C */
  order_type?: string;
}

export interface OrdersOrder {
  code?: string;
  created_at?: string;
  creator_id?: string;
  customer_id?: number;
  deleted_at?: string;
  final_amount?: number;
  id?: string;
  items?: OrdersOrderItem[];
  note?: string;
  order_type?: string;
  payment_status?: string;
  status?: string;
  total_amount?: number;
  updated_at?: string;
}

export interface OrdersOrderItem {
  base_quantity?: number;
  batch_no?: string;
  conversion_factor?: number;
  created_at?: string;
  deleted_at?: string;
  discount?: number;
  expiry_date?: string;
  id?: string;
  is_gift?: boolean;
  note?: string;
  order_id?: string;
  product_id?: number;
  quantity?: number;
  quantity_picked?: number;
  quantity_returned?: number;
  total_line?: number;
  unit_price?: number;
  /** Unit of Measure */
  uom?: string;
}

export interface OrdersUpdateOrderRequest {
  note?: string;
  payment_status?: string;
  status?: string;
}

export interface PayrollsCreatePayrollRequest {
  allowances?: number;
  base_salary: number;
  bonuses?: number;
  deductions?: number;
  employee_id: string;
  insurance_amount?: number;
  net_salary: number;
  overtime_pay?: number;
  period_month: number;
  period_year: number;
  tax_amount?: number;
}

export interface PayrollsPayroll {
  allowances?: number;
  base_salary?: number;
  bonuses?: number;
  created_at?: string;
  deductions?: number;
  employee_id?: string;
  id?: string;
  insurance_amount?: number;
  net_salary?: number;
  overtime_pay?: number;
  payment_date?: string;
  period_month?: number;
  period_year?: number;
  status?: string;
  tax_amount?: number;
  updated_at?: string;
}

export interface PayrollsUpdatePayrollRequest {
  allowances?: number;
  base_salary?: number;
  bonuses?: number;
  deductions?: number;
  insurance_amount?: number;
  net_salary?: number;
  overtime_pay?: number;
  payment_date?: string;
  status?: string;
  tax_amount?: number;
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

export interface PromotionsCreatePromotionRequest {
  code: string;
  end_date: string;
  name: string;
  rules: string[];
  start_date: string;
}

export interface PromotionsPromotion {
  code?: string;
  created_at?: string;
  end_date?: string;
  id?: string;
  name?: string;
  rules?: string[];
  start_date?: string;
  status?: string;
}

export interface PromotionsUpdatePromotionRequest {
  end_date?: string;
  name?: string;
  rules?: string[];
  start_date?: string;
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

export interface ShippingPartnersCreateShippingPartnerRequest {
  api_config?: string[];
  code: string;
  name: string;
  partner_type: string;
  tracking_url_template?: string;
}

export interface ShippingPartnersShippingPartner {
  api_config?: string[];
  code?: string;
  created_at?: string;
  deleted_at?: string;
  id?: string;
  name?: string;
  partner_type?: string;
  status?: string;
  tracking_url_template?: string;
  updated_at?: string;
}

export interface ShippingPartnersUpdateShippingPartnerRequest {
  api_config?: string[];
  code?: string;
  name?: string;
  partner_type?: string;
  status?: string;
  tracking_url_template?: string;
}

export interface SystemConfigsCreateSystemConfigRequest {
  config_key: string;
  config_value: string[];
  description?: string;
  updated_by?: string;
}

export interface SystemConfigsSystemConfig {
  config_key?: string;
  config_value?: string[];
  description?: string;
  id?: string;
  updated_at?: string;
  updated_by?: string;
}

export interface SystemConfigsUpdateSystemConfigRequest {
  config_value?: string[];
  description?: string;
  updated_by?: string;
}

export interface TimeAttendanceCreateTimeAttendanceRequest {
  check_in?: string;
  check_out?: string;
  date: string;
  device_info?: string[];
  employee_id: string;
  location?: string[];
  note?: string;
  overtime_hours?: number;
  shift_type?: string;
  status?: string;
}

export interface TimeAttendanceTimeAttendance {
  check_in?: string;
  check_out?: string;
  created_at?: string;
  date?: string;
  device_info?: string[];
  employee_id?: string;
  id?: string;
  location?: string[];
  note?: string;
  overtime_hours?: number;
  shift_type?: string;
  status?: string;
  updated_at?: string;
}

export interface TimeAttendanceUpdateTimeAttendanceRequest {
  check_in?: string;
  check_out?: string;
  device_info?: string[];
  location?: string[];
  note?: string;
  overtime_hours?: number;
  shift_type?: string;
  status?: string;
}

export interface TrainingCoursesCreateTrainingCourseRequest {
  content_type: string;
  content_url?: string;
  passing_score?: number;
  title: string;
}

export interface TrainingCoursesTrainingCourse {
  content_type?: string;
  content_url?: string;
  created_at?: string;
  deleted_at?: GormDeletedAt;
  id?: number;
  passing_score?: number;
  status?: string;
  title?: string;
  updated_at?: string;
}

export interface TrainingCoursesUpdateTrainingCourseRequest {
  content_type?: string;
  content_url?: string;
  passing_score?: number;
  status?: string;
  title?: string;
}

export interface UserNotificationsCreateFCMTokenRequest {
  device_id?: string;
  device_type?: string;
  fcm_token: string;
  user_id: string;
}

export interface UserNotificationsCreateSocialMappingRequest {
  social_avatar?: string;
  social_id: string;
  social_name?: string;
  social_provider: string;
  user_id: string;
}

export interface UserNotificationsUpdateFCMTokenRequest {
  device_id?: string;
  device_type?: string;
  fcm_token?: string;
}

export interface UserNotificationsUserFCMToken {
  created_at?: string;
  device_id?: string;
  device_type?: string;
  fcm_token?: string;
  id?: number;
  updated_at?: string;
  user_id?: string;
}

export interface UserNotificationsUserSocialMapping {
  created_at?: string;
  id?: number;
  social_avatar?: string;
  social_id?: string;
  social_name?: string;
  social_provider?: string;
  user_id?: string;
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
  manager?: string;
  name: string;
  status?: string;
  type?: string;
}

export interface WarehousesUpdateWarehouseRequest {
  address?: string;
  manager?: string;
  name?: string;
  status?: string;
  type?: string;
}

export interface WarehousesWarehouse {
  address?: string;
  created_at?: string;
  deleted_at?: string;
  id?: number;
  manager?: string;
  name?: string;
  status?: string;
  type?: string;
  updated_at?: string;
}

export interface WorkShiftsCreateShiftAssignmentRequest {
  is_overtime?: boolean;
  shift_id: number;
  status?: string;
  user_id: string;
  work_date: string;
}

export interface WorkShiftsCreateShiftHandoverRequest {
  actual_cash_submitted: number;
  assignment_id?: number;
  branch_id: number;
  system_cash_amount?: number;
  system_cod_amount?: number;
  user_id: string;
}

export interface WorkShiftsCreateWorkShiftRequest {
  branch_id: number;
  end_time: string;
  is_active?: boolean;
  name: string;
  start_time: string;
}

export interface WorkShiftsShiftAssignment {
  created_at?: string;
  id?: number;
  is_overtime?: boolean;
  shift_id?: number;
  status?: string;
  user_id?: string;
  work_date?: string;
}

export interface WorkShiftsShiftHandover {
  actual_cash_submitted?: number;
  assignment_id?: number;
  branch_id?: number;
  created_at?: string;
  finance_transaction_id?: number;
  id?: string;
  status?: string;
  system_cash_amount?: number;
  system_cod_amount?: number;
  user_id?: string;
}

export interface WorkShiftsUpdateShiftAssignmentRequest {
  is_overtime?: boolean;
  status?: string;
}

export interface WorkShiftsUpdateShiftHandoverRequest {
  finance_transaction_id?: number;
  status?: string;
}

export interface WorkShiftsUpdateWorkShiftRequest {
  branch_id?: number;
  end_time?: string;
  is_active?: boolean;
  name?: string;
  start_time?: string;
}

export interface WorkShiftsWorkShift {
  branch_id?: number;
  created_at?: string;
  end_time?: string;
  id?: number;
  is_active?: boolean;
  name?: string;
  start_time?: string;
}

// Auto-generated flattened exports for openapi-typescript v5 backwards compatibility

// Aliases just in case
export type Category = components['definitions']['categories.Category'];
export type Batch = components['definitions']['batches.Batch'];
export type Product = components['definitions']['products.Product'];
