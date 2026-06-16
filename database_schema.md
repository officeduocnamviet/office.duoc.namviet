# Database Schema

### Table: accounting_journals
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **entry_date**: date (Nullable: NO)
- **doc_type**: text (Nullable: NO)
- **source_ref_id**: text (Nullable: YES)
- **description**: text (Nullable: YES)
- **account_debit**: text (Nullable: NO)
- **account_credit**: text (Nullable: NO)
- **amount**: numeric (Nullable: NO)
- **posted_by**: uuid (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: agent_workflows
- **intent_code**: text (Nullable: NO)
- **description**: text (Nullable: NO)
- **required_permission**: text (Nullable: NO)
- **draft_only**: boolean (Nullable: YES) DEFAULT true
- **api_endpoint**: text (Nullable: NO)
- **is_active**: boolean (Nullable: YES) DEFAULT true

### Table: ai_agent_memories
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **user_id**: uuid (Nullable: YES)
- **customer_id**: bigint (Nullable: YES)
- **memory_type**: text (Nullable: NO)
- **memory_text**: text (Nullable: NO)
- **embedding**: USER-DEFINED (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: appointments
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **customer_id**: bigint (Nullable: NO)
- **doctor_id**: uuid (Nullable: YES)
- **room_id**: bigint (Nullable: YES)
- **service_type**: text (Nullable: NO)
- **appointment_time**: timestamp with time zone (Nullable: NO)
- **check_in_time**: timestamp with time zone (Nullable: YES)
- **status**: text (Nullable: NO) DEFAULT 'pending'::text
- **symptoms**: jsonb (Nullable: YES) DEFAULT '[]'::jsonb
- **note**: text (Nullable: YES)
- **created_by**: uuid (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **deleted_at**: timestamp with time zone (Nullable: YES)

### Table: approval_requests
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **requester_id**: uuid (Nullable: NO)
- **module**: text (Nullable: NO)
- **reference_id**: text (Nullable: YES)
- **payload**: jsonb (Nullable: NO)
- **status**: text (Nullable: NO) DEFAULT 'pending'::text
- **current_step**: integer (Nullable: YES) DEFAULT 1
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: approval_steps
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **request_id**: uuid (Nullable: NO)
- **step_order**: integer (Nullable: NO)
- **approver_id**: uuid (Nullable: YES)
- **approver_role_id**: uuid (Nullable: YES)
- **status**: text (Nullable: NO) DEFAULT 'pending'::text
- **note**: text (Nullable: YES)
- **processed_at**: timestamp with time zone (Nullable: YES)

### Table: assets
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **name**: text (Nullable: NO)
- **category**: text (Nullable: NO)
- **purchase_price**: numeric (Nullable: NO)
- **purchase_date**: date (Nullable: NO)
- **depreciation_months**: integer (Nullable: NO)
- **status**: text (Nullable: NO) DEFAULT 'active'::text
- **assigned_to**: uuid (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: attendance_logs
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **user_id**: uuid (Nullable: NO)
- **branch_id**: bigint (Nullable: YES)
- **check_in_time**: timestamp with time zone (Nullable: NO) DEFAULT now()
- **check_in_ip**: text (Nullable: YES)
- **check_in_lat**: numeric (Nullable: YES)
- **check_in_lng**: numeric (Nullable: YES)
- **check_out_time**: timestamp with time zone (Nullable: YES)
- **check_out_ip**: text (Nullable: YES)
- **check_out_lat**: numeric (Nullable: YES)
- **check_out_lng**: numeric (Nullable: YES)
- **is_valid**: boolean (Nullable: YES) DEFAULT false
- **working_hours**: numeric (Nullable: YES) DEFAULT 0
- **note**: text (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: batches
- **id**: bigint (Nullable: NO)
- **product_id**: bigint (Nullable: NO)
- **batch_code**: text (Nullable: NO)
- **expiry_date**: date (Nullable: NO)
- **manufacturing_date**: date (Nullable: YES)
- **inbound_price**: numeric (Nullable: YES) DEFAULT 0
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **deleted_at**: timestamp with time zone (Nullable: YES)

### Table: categories
- **id**: bigint (Nullable: NO)
- **name**: text (Nullable: NO)
- **slug**: text (Nullable: NO)
- **parent_id**: bigint (Nullable: YES)
- **status**: text (Nullable: NO) DEFAULT 'active'::text
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **deleted_at**: timestamp with time zone (Nullable: YES)

### Table: chart_of_accounts
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **account_code**: text (Nullable: NO)
- **name**: text (Nullable: NO)
- **parent_id**: uuid (Nullable: YES)
- **type**: text (Nullable: NO)
- **balance_type**: text (Nullable: NO)
- **status**: text (Nullable: NO) DEFAULT 'active'::text
- **allow_posting**: boolean (Nullable: NO) DEFAULT true
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **deleted_at**: timestamp with time zone (Nullable: YES)

### Table: chat_messages
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **session_id**: uuid (Nullable: NO)
- **role**: text (Nullable: NO)
- **content**: text (Nullable: YES)
- **intent**: text (Nullable: YES)
- **entities**: jsonb (Nullable: YES)
- **llm_meta**: jsonb (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: chat_sessions
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **user_id**: uuid (Nullable: NO)
- **platform**: text (Nullable: NO) DEFAULT 'web'::text
- **status**: text (Nullable: NO) DEFAULT 'bot'::text
- **context**: jsonb (Nullable: NO) DEFAULT '{}'::jsonb
- **assigned_sales_id**: uuid (Nullable: YES)
- **started_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **last_activity_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **closed_at**: timestamp with time zone (Nullable: YES)

### Table: clinical_queues
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **appointment_id**: uuid (Nullable: YES)
- **customer_id**: bigint (Nullable: NO)
- **doctor_id**: uuid (Nullable: YES)
- **queue_number**: integer (Nullable: NO)
- **status**: text (Nullable: NO) DEFAULT 'waiting'::text
- **priority_level**: text (Nullable: NO) DEFAULT 'normal'::text
- **checked_in_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **deleted_at**: timestamp with time zone (Nullable: YES)

### Table: companies
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **tax_code**: text (Nullable: NO)
- **name**: text (Nullable: NO)
- **short_name**: text (Nullable: YES)
- **address**: text (Nullable: NO)
- **phone**: text (Nullable: NO)
- **email**: text (Nullable: NO)
- **logo_url**: text (Nullable: YES)
- **representative_name**: text (Nullable: NO)
- **business_license_url**: ARRAY (Nullable: YES)
- **mission**: text (Nullable: YES)
- **vision**: text (Nullable: YES)
- **status**: text (Nullable: YES) DEFAULT 'active'::text
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **deleted_at**: timestamp with time zone (Nullable: YES)

### Table: customer_vaccination_records
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **customer_id**: bigint (Nullable: NO)
- **appointment_id**: uuid (Nullable: YES)
- **product_id**: bigint (Nullable: NO)
- **dose_number**: integer (Nullable: NO) DEFAULT 1
- **expected_date**: date (Nullable: NO)
- **actual_date**: date (Nullable: YES)
- **status**: text (Nullable: NO) DEFAULT 'pending'::text
- **consulted_by**: uuid (Nullable: YES)
- **administered_by**: uuid (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **deleted_at**: timestamp with time zone (Nullable: YES)

### Table: customer_vouchers
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **promotion_id**: uuid (Nullable: YES)
- **customer_id**: bigint (Nullable: YES)
- **is_used**: boolean (Nullable: YES) DEFAULT false
- **used_at**: timestamp with time zone (Nullable: YES)
- **order_id**: uuid (Nullable: YES)

### Table: customers
- **id**: bigint (Nullable: NO)
- **customer_code**: text (Nullable: YES)
- **name**: text (Nullable: NO)
- **customer_type**: text (Nullable: NO) DEFAULT 'B2C'::text
- **phone**: text (Nullable: YES)
- **email**: text (Nullable: YES)
- **address**: text (Nullable: YES)
- **status**: text (Nullable: NO) DEFAULT 'active'::text
- **dob**: date (Nullable: YES)
- **gender**: text (Nullable: YES)
- **cccd**: text (Nullable: YES)
- **loyalty_points**: integer (Nullable: YES) DEFAULT 0
- **b2b_metadata**: jsonb (Nullable: YES) DEFAULT '{}'::jsonb
- **current_debt**: numeric (Nullable: YES) DEFAULT 0
- **updated_by**: uuid (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **deleted_at**: timestamp with time zone (Nullable: YES)

### Table: employment_contracts
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **user_id**: uuid (Nullable: NO)
- **contract_code**: text (Nullable: NO)
- **base_salary**: numeric (Nullable: NO) DEFAULT 0
- **standard_working_days**: integer (Nullable: NO) DEFAULT 26
- **kpi_conversion_rate**: numeric (Nullable: YES) DEFAULT 0
- **commission_rate_percent**: numeric (Nullable: YES) DEFAULT 0
- **tax_deduction_amount**: numeric (Nullable: YES) DEFAULT 0
- **insurance_deduction_amount**: numeric (Nullable: YES) DEFAULT 0
- **valid_from**: date (Nullable: NO)
- **valid_to**: date (Nullable: YES)
- **status**: text (Nullable: NO) DEFAULT 'active'::text
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: finance_transactions
- **id**: bigint (Nullable: NO)
- **code**: text (Nullable: NO)
- **transaction_date**: timestamp with time zone (Nullable: NO) DEFAULT now()
- **flow**: text (Nullable: NO)
- **business_type**: text (Nullable: NO) DEFAULT 'other'::text
- **category_id**: bigint (Nullable: YES)
- **amount**: numeric (Nullable: NO)
- **fund_account_id**: bigint (Nullable: NO)
- **partner_type**: text (Nullable: YES)
- **partner_id**: text (Nullable: YES)
- **partner_name_cache**: text (Nullable: YES)
- **ref_type**: text (Nullable: YES)
- **ref_id**: text (Nullable: YES)
- **description**: text (Nullable: YES)
- **evidence_url**: text (Nullable: YES)
- **status**: text (Nullable: NO) DEFAULT 'pending'::text
- **cash_tally**: jsonb (Nullable: YES)
- **ref_advance_id**: bigint (Nullable: YES)
- **target_bank_info**: jsonb (Nullable: YES)
- **bank_reference_id**: text (Nullable: YES)
- **book_type**: text (Nullable: NO) DEFAULT 'BOTH'::text
- **is_posted**: boolean (Nullable: NO) DEFAULT false
- **created_by**: uuid (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **deleted_at**: timestamp with time zone (Nullable: YES)

### Table: fund_accounts
- **id**: bigint (Nullable: NO)
- **name**: text (Nullable: NO)
- **type**: text (Nullable: NO)
- **location**: text (Nullable: YES)
- **account_number**: text (Nullable: YES)
- **bank_id**: bigint (Nullable: YES)
- **initial_balance**: numeric (Nullable: NO) DEFAULT 0
- **balance**: numeric (Nullable: NO) DEFAULT 0
- **currency**: text (Nullable: YES) DEFAULT 'VND'::text
- **status**: text (Nullable: NO) DEFAULT 'active'::text
- **bank_info**: jsonb (Nullable: YES)
- **description**: text (Nullable: YES)
- **account_id**: text (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **deleted_at**: timestamp with time zone (Nullable: YES)

### Table: internal_channels
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **name**: text (Nullable: YES)
- **type**: text (Nullable: NO)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: internal_messages
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **channel_id**: uuid (Nullable: NO)
- **sender_id**: uuid (Nullable: YES)
- **content**: text (Nullable: YES)
- **attachments**: jsonb (Nullable: YES) DEFAULT '[]'::jsonb
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: inventory_batches
- **id**: bigint (Nullable: NO)
- **warehouse_id**: bigint (Nullable: NO)
- **product_id**: bigint (Nullable: NO)
- **batch_id**: bigint (Nullable: NO)
- **quantity**: numeric (Nullable: NO) DEFAULT 0
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: inventory_transactions
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **warehouse_id**: bigint (Nullable: NO)
- **product_id**: bigint (Nullable: NO)
- **batch_id**: bigint (Nullable: YES)
- **type**: text (Nullable: NO)
- **action_group**: text (Nullable: YES)
- **quantity**: numeric (Nullable: NO)
- **unit_price**: numeric (Nullable: YES) DEFAULT 0
- **ref_id**: text (Nullable: YES)
- **description**: text (Nullable: YES)
- **partner_id**: bigint (Nullable: YES)
- **created_by**: uuid (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: manufacturers
- **id**: bigint (Nullable: NO)
- **name**: text (Nullable: NO)
- **slug**: text (Nullable: NO)
- **country**: text (Nullable: YES)
- **logo_url**: text (Nullable: YES)
- **status**: text (Nullable: YES) DEFAULT 'active'::text
- **created_at**: timestamp with time zone (Nullable: NO) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: NO) DEFAULT now()
- **deleted_at**: timestamp with time zone (Nullable: YES)

### Table: marketing_campaigns
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **name**: text (Nullable: NO)
- **objective**: text (Nullable: NO)
- **start_date**: timestamp with time zone (Nullable: NO)
- **end_date**: timestamp with time zone (Nullable: NO)
- **budget**: numeric (Nullable: YES) DEFAULT 0
- **actual_cost**: numeric (Nullable: YES) DEFAULT 0
- **status**: text (Nullable: NO) DEFAULT 'draft'::text
- **target_segment_id**: bigint (Nullable: YES)
- **created_by**: uuid (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **deleted_at**: timestamp with time zone (Nullable: YES)

### Table: medical_knowledge_vectors
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **title**: text (Nullable: NO)
- **content**: text (Nullable: NO)
- **metadata**: jsonb (Nullable: YES) DEFAULT '{}'::jsonb
- **embedding**: USER-DEFINED (Nullable: YES)
- **created_by**: uuid (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: medical_visits
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **appointment_id**: uuid (Nullable: YES)
- **customer_id**: bigint (Nullable: NO)
- **doctor_id**: uuid (Nullable: YES)
- **temperature**: numeric (Nullable: YES)
- **pulse**: integer (Nullable: YES)
- **sp02**: integer (Nullable: YES)
- **bp_systolic**: integer (Nullable: YES)
- **bp_diastolic**: integer (Nullable: YES)
- **weight**: numeric (Nullable: YES)
- **height**: numeric (Nullable: YES)
- **symptoms**: text (Nullable: YES)
- **examination_summary**: text (Nullable: YES)
- **diagnosis**: text (Nullable: YES)
- **icd_code**: text (Nullable: YES)
- **doctor_notes**: text (Nullable: YES)
- **red_flags**: jsonb (Nullable: YES) DEFAULT '[]'::jsonb
- **status**: text (Nullable: NO) DEFAULT 'in_progress'::text
- **created_by**: uuid (Nullable: YES)
- **updated_by**: uuid (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **deleted_at**: timestamp with time zone (Nullable: YES)

### Table: order_items
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **order_id**: uuid (Nullable: NO)
- **product_id**: bigint (Nullable: NO)
- **quantity**: integer (Nullable: NO)
- **uom**: text (Nullable: NO)
- **conversion_factor**: integer (Nullable: YES)
- **base_quantity**: integer (Nullable: YES)
- **unit_price**: numeric (Nullable: NO)
- **discount**: numeric (Nullable: YES) DEFAULT 0
- **is_gift**: boolean (Nullable: YES) DEFAULT false
- **note**: text (Nullable: YES)
- **batch_no**: text (Nullable: YES)
- **expiry_date**: date (Nullable: YES)
- **total_line**: numeric (Nullable: YES)
- **quantity_picked**: integer (Nullable: YES) DEFAULT 0
- **quantity_returned**: integer (Nullable: YES) DEFAULT 0
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **deleted_at**: timestamp with time zone (Nullable: YES)

### Table: orders
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **code**: text (Nullable: NO)
- **customer_id**: bigint (Nullable: YES)
- **creator_id**: uuid (Nullable: YES)
- **status**: text (Nullable: NO) DEFAULT 'PENDING'::text
- **order_type**: text (Nullable: NO) DEFAULT 'B2C'::text
- **total_amount**: numeric (Nullable: YES) DEFAULT 0
- **final_amount**: numeric (Nullable: YES) DEFAULT 0
- **payment_status**: text (Nullable: YES) DEFAULT 'unpaid'::text
- **note**: text (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **deleted_at**: timestamp with time zone (Nullable: YES)

### Table: payroll_items
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **payroll_id**: uuid (Nullable: NO)
- **employee_id**: uuid (Nullable: NO)
- **basic_salary**: numeric (Nullable: NO) DEFAULT 0
- **working_days**: numeric (Nullable: NO) DEFAULT 0
- **kpi_score**: numeric (Nullable: YES) DEFAULT 0
- **bonuses**: numeric (Nullable: YES) DEFAULT 0
- **deductions**: numeric (Nullable: YES) DEFAULT 0
- **net_pay**: numeric (Nullable: NO) DEFAULT 0

### Table: payroll_items_v2
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **payroll_id**: uuid (Nullable: NO)
- **employee_id**: uuid (Nullable: NO)
- **base_salary**: numeric (Nullable: NO)
- **standard_days**: integer (Nullable: NO)
- **actual_days**: numeric (Nullable: NO)
- **prorated_salary**: numeric (Nullable: NO)
- **total_kpi_points**: numeric (Nullable: YES) DEFAULT 0
- **kpi_bonus_amount**: numeric (Nullable: YES) DEFAULT 0
- **total_sales_amount**: numeric (Nullable: YES) DEFAULT 0
- **commission_bonus_amount**: numeric (Nullable: YES) DEFAULT 0
- **other_bonus_amount**: numeric (Nullable: YES) DEFAULT 0
- **tax_deduction**: numeric (Nullable: YES) DEFAULT 0
- **insurance_deduction**: numeric (Nullable: YES) DEFAULT 0
- **net_pay**: numeric (Nullable: NO)
- **employee_agreed**: boolean (Nullable: YES) DEFAULT false
- **employee_note**: text (Nullable: YES)
- **accountant_verified**: boolean (Nullable: YES) DEFAULT false
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: payrolls
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **code**: text (Nullable: NO)
- **month**: integer (Nullable: NO)
- **year**: integer (Nullable: NO)
- **total_basic_salary**: numeric (Nullable: YES) DEFAULT 0
- **total_allowance**: numeric (Nullable: YES) DEFAULT 0
- **total_bonus_kpi**: numeric (Nullable: YES) DEFAULT 0
- **total_deduction**: numeric (Nullable: YES) DEFAULT 0
- **net_pay**: numeric (Nullable: YES) DEFAULT 0
- **status**: text (Nullable: NO) DEFAULT 'draft'::text
- **created_by**: uuid (Nullable: YES)
- **approved_by**: uuid (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: product_inventory
- **id**: bigint (Nullable: NO)
- **product_id**: bigint (Nullable: YES)
- **warehouse_id**: bigint (Nullable: YES)
- **stock_quantity**: numeric (Nullable: NO) DEFAULT 0
- **min_stock**: integer (Nullable: YES) DEFAULT 0
- **max_stock**: integer (Nullable: YES) DEFAULT 0
- **shelf_location**: text (Nullable: YES) DEFAULT 'Chưa xếp'::text
- **location_cabinet**: text (Nullable: YES)
- **location_row**: text (Nullable: YES)
- **location_slot**: text (Nullable: YES)
- **updated_by**: uuid (Nullable: YES)
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: product_units
- **id**: bigint (Nullable: NO)
- **product_id**: bigint (Nullable: YES)
- **unit_name**: text (Nullable: NO)
- **conversion_rate**: integer (Nullable: YES) DEFAULT 1
- **barcode**: text (Nullable: YES)
- **is_base**: boolean (Nullable: YES) DEFAULT false
- **is_direct_sale**: boolean (Nullable: YES) DEFAULT true
- **price_cost**: numeric (Nullable: YES) DEFAULT 0
- **price_sell**: numeric (Nullable: YES) DEFAULT 0
- **unit_type**: text (Nullable: YES) DEFAULT 'retail'::text
- **price**: numeric (Nullable: YES) DEFAULT 0
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **deleted_at**: timestamp with time zone (Nullable: YES)

### Table: product_vectors
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **product_id**: bigint (Nullable: NO)
- **semantic_text**: text (Nullable: NO)
- **embedding**: USER-DEFINED (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: products
- **id**: bigint (Nullable: NO)
- **name**: text (Nullable: NO)
- **sku**: text (Nullable: YES)
- **barcode**: text (Nullable: YES)
- **description**: text (Nullable: YES)
- **active_ingredient**: text (Nullable: YES)
- **image_url**: text (Nullable: YES)
- **status**: text (Nullable: NO) DEFAULT 'active'::text
- **fts**: tsvector (Nullable: YES)
- **category_id**: bigint (Nullable: YES)
- **manufacturer_id**: bigint (Nullable: YES)
- **category_name**: text (Nullable: YES)
- **manufacturer_name**: text (Nullable: YES)
- **distributor_id**: bigint (Nullable: YES)
- **invoice_price**: numeric (Nullable: YES) DEFAULT 0
- **actual_cost**: numeric (Nullable: NO) DEFAULT 0
- **wholesale_unit**: text (Nullable: YES) DEFAULT 'Hộp'::text
- **retail_unit**: text (Nullable: YES) DEFAULT 'Vỉ'::text
- **conversion_factor**: integer (Nullable: YES) DEFAULT 1
- **wholesale_margin_value**: numeric (Nullable: YES) DEFAULT 0
- **wholesale_margin_type**: text (Nullable: YES) DEFAULT '%'::text
- **retail_margin_value**: numeric (Nullable: YES) DEFAULT 0
- **retail_margin_type**: text (Nullable: YES) DEFAULT '%'::text
- **items_per_carton**: integer (Nullable: YES) DEFAULT 1
- **carton_weight**: numeric (Nullable: YES) DEFAULT 0
- **carton_dimensions**: text (Nullable: YES)
- **purchasing_policy**: text (Nullable: YES) DEFAULT 'ALLOW_LOOSE'::text
- **registration_number**: text (Nullable: YES)
- **packing_spec**: text (Nullable: YES)
- **stock_management_type**: text (Nullable: YES) DEFAULT 'lot_date'::text
- **wholesale_margin_rate**: numeric (Nullable: YES) DEFAULT 0
- **retail_margin_rate**: numeric (Nullable: YES) DEFAULT 0
- **usage_instructions**: jsonb (Nullable: YES) DEFAULT '{}'::jsonb
- **stock_status**: text (Nullable: YES) DEFAULT 'in_stock'::text
- **product_images**: ARRAY (Nullable: YES) DEFAULT '{}'::text[]
- **updated_by**: uuid (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **deleted_at**: timestamp with time zone (Nullable: YES)

### Table: promotions
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **code**: text (Nullable: NO)
- **name**: text (Nullable: NO)
- **rules**: jsonb (Nullable: NO)
- **start_date**: timestamp with time zone (Nullable: NO)
- **end_date**: timestamp with time zone (Nullable: NO)
- **status**: text (Nullable: YES) DEFAULT 'active'::text
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: roles
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **name**: text (Nullable: NO)
- **description**: text (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **permissions**: jsonb (Nullable: YES) DEFAULT '[]'::jsonb

### Table: shift_assignments
- **id**: bigint (Nullable: NO)
- **shift_id**: bigint (Nullable: NO)
- **user_id**: uuid (Nullable: NO)
- **work_date**: date (Nullable: NO)
- **status**: text (Nullable: NO) DEFAULT 'scheduled'::text
- **is_overtime**: boolean (Nullable: YES) DEFAULT false
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: shift_handovers
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **assignment_id**: bigint (Nullable: YES)
- **user_id**: uuid (Nullable: NO)
- **branch_id**: bigint (Nullable: NO)
- **system_cash_amount**: numeric (Nullable: NO) DEFAULT 0
- **system_cod_amount**: numeric (Nullable: NO) DEFAULT 0
- **actual_cash_submitted**: numeric (Nullable: NO)
- **status**: text (Nullable: NO) DEFAULT 'pending_finance'::text
- **finance_transaction_id**: bigint (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: system_audit_logs
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **table_name**: text (Nullable: NO)
- **record_id**: text (Nullable: NO)
- **action**: text (Nullable: NO)
- **old_data**: jsonb (Nullable: YES)
- **new_data**: jsonb (Nullable: YES)
- **performed_by**: uuid (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: NO) DEFAULT now()

### Table: system_audit_logs_2026_06
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **table_name**: text (Nullable: NO)
- **record_id**: text (Nullable: NO)
- **action**: text (Nullable: NO)
- **old_data**: jsonb (Nullable: YES)
- **new_data**: jsonb (Nullable: YES)
- **performed_by**: uuid (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: NO) DEFAULT now()

### Table: system_audit_logs_2026_07
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **table_name**: text (Nullable: NO)
- **record_id**: text (Nullable: NO)
- **action**: text (Nullable: NO)
- **old_data**: jsonb (Nullable: YES)
- **new_data**: jsonb (Nullable: YES)
- **performed_by**: uuid (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: NO) DEFAULT now()

### Table: system_configs
- **key**: text (Nullable: NO)
- **value**: jsonb (Nullable: NO)
- **description**: text (Nullable: YES)
- **updated_by**: uuid (Nullable: YES)
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: third_party_connections
- **id**: bigint (Nullable: NO)
- **provider**: text (Nullable: NO)
- **access_token**: text (Nullable: YES)
- **refresh_token**: text (Nullable: YES)
- **config**: jsonb (Nullable: YES) DEFAULT '{}'::jsonb
- **status**: text (Nullable: YES) DEFAULT 'connected'::text
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: training_courses
- **id**: bigint (Nullable: NO)
- **title**: text (Nullable: NO)
- **content_type**: text (Nullable: NO)
- **content_url**: text (Nullable: YES)
- **passing_score**: integer (Nullable: YES)
- **status**: text (Nullable: NO) DEFAULT 'active'::text
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **deleted_at**: timestamp with time zone (Nullable: YES)

### Table: user_fcm_tokens
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **user_id**: uuid (Nullable: NO)
- **token**: text (Nullable: NO)
- **device_info**: text (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: user_social_mappings
- **user_id**: uuid (Nullable: NO)
- **telegram_id**: text (Nullable: YES)
- **zalo_id**: text (Nullable: YES)
- **is_verified**: boolean (Nullable: YES) DEFAULT false
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: users
- **id**: uuid (Nullable: NO)
- **email**: text (Nullable: YES)
- **full_name**: text (Nullable: YES)
- **avatar_url**: text (Nullable: YES)
- **employee_code**: text (Nullable: YES)
- **position**: text (Nullable: YES)
- **status**: text (Nullable: NO) DEFAULT 'pending_approval'::text
- **dob**: date (Nullable: YES)
- **phone**: text (Nullable: YES)
- **gender**: text (Nullable: YES)
- **cccd**: text (Nullable: YES)
- **cccd_issue_date**: date (Nullable: YES)
- **address**: text (Nullable: YES)
- **marital_status**: text (Nullable: YES)
- **cccd_front_url**: text (Nullable: YES)
- **cccd_back_url**: text (Nullable: YES)
- **education_level**: text (Nullable: YES)
- **specialization**: text (Nullable: YES)
- **bank_name**: text (Nullable: YES)
- **bank_account_number**: text (Nullable: YES)
- **bank_account_name**: text (Nullable: YES)
- **hobbies**: text (Nullable: YES)
- **limitations**: text (Nullable: YES)
- **strengths**: text (Nullable: YES)
- **needs**: text (Nullable: YES)
- **work_state**: text (Nullable: YES) DEFAULT 'working'::text
- **role_id**: uuid (Nullable: YES)
- **company_id**: uuid (Nullable: YES)
- **warehouse_id**: bigint (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **profile_updated_at**: timestamp with time zone (Nullable: YES)
- **deleted_at**: timestamp with time zone (Nullable: YES)

### Table: vendor_product_mappings
- **id**: bigint (Nullable: NO)
- **vendor_tax_code**: text (Nullable: NO)
- **vendor_product_name**: text (Nullable: NO)
- **vendor_unit**: text (Nullable: YES)
- **internal_product_id**: bigint (Nullable: NO)
- **internal_unit**: text (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **last_used_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **updated_by**: uuid (Nullable: YES)

### Table: warehouses
- **id**: bigint (Nullable: NO)
- **key**: text (Nullable: NO)
- **name**: text (Nullable: NO)
- **unit**: text (Nullable: NO) DEFAULT 'Hộp'::text
- **address**: text (Nullable: YES)
- **type**: text (Nullable: NO) DEFAULT 'retail'::text
- **latitude**: numeric (Nullable: YES)
- **longitude**: numeric (Nullable: YES)
- **code**: text (Nullable: YES)
- **manager**: text (Nullable: YES)
- **phone**: text (Nullable: YES)
- **status**: text (Nullable: NO) DEFAULT 'active'::text
- **company_id**: uuid (Nullable: YES)
- **outlet_type**: text (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()
- **deleted_at**: timestamp with time zone (Nullable: YES)

### Table: webhook_logs
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **provider**: text (Nullable: NO)
- **payload**: jsonb (Nullable: NO)
- **status**: text (Nullable: NO) DEFAULT 'processing'::text
- **error_message**: text (Nullable: YES)
- **processed_at**: timestamp with time zone (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

### Table: work_shifts
- **id**: bigint (Nullable: NO)
- **branch_id**: bigint (Nullable: NO)
- **name**: text (Nullable: NO)
- **start_time**: time without time zone (Nullable: NO)
- **end_time**: time without time zone (Nullable: NO)
- **is_active**: boolean (Nullable: YES) DEFAULT true
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now()

