# encoding: UTF-8
# This file is auto-generated from the current state of the database. Instead
# of editing this file, please use the migrations feature of Active Record to
# incrementally modify your database, and then regenerate this schema definition.
#
# Note that this schema.rb definition is the authoritative source for your
# database schema. If you need to create the application database on another
# system, you should be using db:schema:load, not running all the migrations
# from scratch. The latter is a flawed and unsustainable approach (the more migrations
# you'll amass, the slower it'll run and the greater likelihood for issues).
#
# It's strongly recommended to check this file into your version control system.

ActiveRecord::Schema.define(:version => 20121029054455) do

  create_table "attachments", :force => true do |t|
    t.string   "uuid",                  :limit => 36
    t.string   "document_file_name"
    t.string   "document_content_type"
    t.integer  "document_file_size"
    t.datetime "document_updated_at"
    t.integer  "lock_version",                        :default => 0
    t.string   "original_file_name"
    t.string   "create_user_uuid",      :limit => 36
  end

  add_index "attachments", ["uuid"], :name => "index_attachments_on_uuid"

  create_table "audits", :force => true do |t|
    t.integer  "version"
    t.integer  "lock_version",                   :default => 0
    t.string   "kind",             :limit => 36
    t.datetime "updated_at"
    t.string   "locale"
    t.string   "update_user_uuid"
    t.string   "uuid"
    t.string   "grid_uuid",        :limit => 36
  end

  add_index "audits", ["update_user_uuid"], :name => "index_audits_on_update_user_uuid"
  add_index "audits", ["uuid"], :name => "index_audits_on_uuid"

  create_table "column_locs", :force => true do |t|
    t.string   "uuid",             :limit => 36
    t.integer  "version"
    t.string   "locale",           :limit => 10
    t.string   "base_locale",      :limit => 10
    t.string   "name"
    t.text     "description"
    t.integer  "lock_version",                   :default => 0
    t.datetime "created_at"
    t.datetime "updated_at"
    t.string   "create_user_uuid", :limit => 36
    t.string   "update_user_uuid", :limit => 36
  end

  add_index "column_locs", ["uuid", "version", "locale"], :name => "index_column_locs_on_uuid_and_version_and_locale"

  create_table "column_mappings", :force => true do |t|
    t.string   "uuid",             :limit => 36
    t.date     "begin"
    t.date     "end"
    t.integer  "version"
    t.boolean  "enabled"
    t.integer  "lock_version",                   :default => 0
    t.string   "column_uuid",      :limit => 36
    t.string   "create_user_uuid", :limit => 36
    t.string   "update_user_uuid", :limit => 36
    t.string   "db_column"
    t.datetime "created_at"
    t.datetime "updated_at"
    t.string   "uri"
  end

  add_index "column_mappings", ["column_uuid"], :name => "index_column_mappings_on_column_uuid"
  add_index "column_mappings", ["uri"], :name => "index_column_mappings_on_uri"

  create_table "columns", :force => true do |t|
    t.date     "begin"
    t.date     "end"
    t.integer  "version"
    t.integer  "number"
    t.integer  "display"
    t.string   "kind"
    t.datetime "created_at"
    t.datetime "updated_at"
    t.integer  "lock_version",                      :default => 0
    t.string   "uuid",                :limit => 36
    t.string   "grid_uuid",           :limit => 36
    t.string   "grid_reference_uuid", :limit => 36
    t.string   "create_user_uuid",    :limit => 36
    t.string   "update_user_uuid",    :limit => 36
    t.boolean  "enabled"
    t.boolean  "required"
    t.string   "regex"
    t.string   "uri"
  end

  add_index "columns", ["grid_uuid"], :name => "index_columns_on_grid_uuid"
  add_index "columns", ["uri"], :name => "index_columns_on_uri"
  add_index "columns", ["uuid", "begin", "end"], :name => "index_columns_on_uuid_and_begin_and_end"
  add_index "columns", ["uuid"], :name => "index_columns_on_uuid"

  create_table "grid_locs", :force => true do |t|
    t.string   "uuid",             :limit => 36
    t.integer  "version"
    t.string   "locale",           :limit => 10
    t.string   "base_locale",      :limit => 10
    t.string   "name"
    t.text     "description"
    t.integer  "lock_version",                   :default => 0
    t.datetime "created_at"
    t.datetime "updated_at"
    t.string   "create_user_uuid", :limit => 36
    t.string   "update_user_uuid", :limit => 36
  end

  add_index "grid_locs", ["uuid", "version", "locale"], :name => "index_grid_locs_on_uuid_and_version_and_locale"

  create_table "grid_mappings", :force => true do |t|
    t.string   "uuid",             :limit => 36
    t.date     "begin"
    t.date     "end"
    t.integer  "version"
    t.boolean  "enabled"
    t.integer  "lock_version",                   :default => 0
    t.string   "grid_uuid",        :limit => 36
    t.string   "create_user_uuid", :limit => 36
    t.string   "update_user_uuid", :limit => 36
    t.string   "db_table"
    t.string   "db_loc_table"
    t.datetime "created_at"
    t.datetime "updated_at"
    t.string   "uri"
  end

  add_index "grid_mappings", ["grid_uuid"], :name => "index_grid_mappings_on_grid_uuid"
  add_index "grid_mappings", ["uri"], :name => "index_grid_mappings_on_uri"

  create_table "grids", :force => true do |t|
    t.date     "begin"
    t.date     "end"
    t.integer  "version"
    t.datetime "created_at"
    t.datetime "updated_at"
    t.integer  "lock_version",                   :default => 0
    t.string   "uuid",             :limit => 36
    t.string   "workspace_uuid",   :limit => 36
    t.string   "create_user_uuid", :limit => 36
    t.string   "update_user_uuid", :limit => 36
    t.boolean  "enabled"
    t.boolean  "has_name",                       :default => true
    t.boolean  "has_description",                :default => true
    t.string   "uri"
    t.string   "template_uuid",    :limit => 36
    t.string   "display_uuid",     :limit => 36
    t.string   "sort_uuid",        :limit => 36
  end

  add_index "grids", ["uri"], :name => "index_grids_on_uri"
  add_index "grids", ["uuid", "begin", "end"], :name => "index_grids_on_uuid_and_begin_and_end"
  add_index "grids", ["uuid"], :name => "index_grids_on_uuid"
  add_index "grids", ["workspace_uuid"], :name => "index_grids_on_workspace_uuid"

  create_table "role_locs", :force => true do |t|
    t.string   "uuid",             :limit => 36
    t.integer  "version"
    t.string   "locale",           :limit => 10
    t.string   "base_locale",      :limit => 10
    t.string   "name"
    t.text     "description"
    t.integer  "lock_version",                   :default => 0
    t.datetime "created_at"
    t.datetime "updated_at"
    t.string   "create_user_uuid", :limit => 36
    t.string   "update_user_uuid", :limit => 36
  end

  add_index "role_locs", ["uuid", "version", "locale"], :name => "index_role_locs_on_uuid_and_version_and_locale"

  create_table "roles", :force => true do |t|
    t.string   "uuid",             :limit => 36
    t.date     "begin"
    t.date     "end"
    t.integer  "version"
    t.boolean  "enabled"
    t.integer  "lock_version",                   :default => 0
    t.string   "create_user_uuid", :limit => 36
    t.string   "update_user_uuid", :limit => 36
    t.datetime "created_at"
    t.datetime "updated_at"
    t.string   "uri"
  end

  add_index "roles", ["uri"], :name => "index_roles_on_uri"
  add_index "roles", ["uuid", "begin", "end"], :name => "index_roles_on_uuid_and_begin_and_end"
  add_index "roles", ["uuid"], :name => "index_roles_on_uuid"

  create_table "row_locs", :force => true do |t|
    t.string   "uuid",             :limit => 36
    t.integer  "version"
    t.string   "locale",           :limit => 10
    t.string   "base_locale",      :limit => 10
    t.string   "name"
    t.text     "description"
    t.integer  "lock_version",                   :default => 0
    t.datetime "created_at"
    t.datetime "updated_at"
    t.string   "create_user_uuid", :limit => 36
    t.string   "update_user_uuid", :limit => 36
  end

  add_index "row_locs", ["uuid", "version", "locale"], :name => "index_row_locs_on_uuid_and_version_and_locale"

  create_table "rows", :force => true do |t|
    t.date     "begin"
    t.date     "end"
    t.integer  "version"
    t.integer  "lock_version"
    t.string   "value1"
    t.string   "value2"
    t.string   "value3"
    t.string   "value4"
    t.string   "value5"
    t.string   "value6"
    t.string   "value7"
    t.string   "value8"
    t.string   "value9"
    t.string   "value10"
    t.string   "value11"
    t.string   "value12"
    t.string   "value13"
    t.string   "value14"
    t.string   "value15"
    t.string   "value16"
    t.string   "value17"
    t.string   "value18"
    t.string   "value19"
    t.string   "value20"
    t.datetime "created_at"
    t.datetime "updated_at"
    t.string   "uuid",             :limit => 36
    t.string   "grid_uuid",        :limit => 36
    t.string   "create_user_uuid", :limit => 36
    t.string   "update_user_uuid", :limit => 36
    t.string   "row_uuid1",        :limit => 36
    t.string   "row_uuid2",        :limit => 36
    t.string   "row_uuid3",        :limit => 36
    t.string   "row_uuid4",        :limit => 36
    t.string   "row_uuid5",        :limit => 36
    t.string   "row_uuid6",        :limit => 36
    t.string   "row_uuid7",        :limit => 36
    t.string   "row_uuid8",        :limit => 36
    t.string   "row_uuid9",        :limit => 36
    t.string   "row_uuid10",       :limit => 36
    t.string   "row_uuid11",       :limit => 36
    t.string   "row_uuid12",       :limit => 36
    t.string   "row_uuid13",       :limit => 36
    t.string   "row_uuid14",       :limit => 36
    t.string   "row_uuid15",       :limit => 36
    t.string   "row_uuid16",       :limit => 36
    t.string   "row_uuid17",       :limit => 36
    t.string   "row_uuid18",       :limit => 36
    t.string   "row_uuid19",       :limit => 36
    t.string   "row_uuid20",       :limit => 36
    t.date     "date1"
    t.date     "date2"
    t.date     "date3"
    t.date     "date4"
    t.date     "date5"
    t.date     "date6"
    t.date     "date7"
    t.date     "date8"
    t.date     "date9"
    t.date     "date10"
    t.date     "date11"
    t.date     "date12"
    t.date     "date13"
    t.date     "date14"
    t.date     "date15"
    t.date     "date16"
    t.date     "date17"
    t.date     "date18"
    t.date     "date19"
    t.date     "date20"
    t.boolean  "enabled"
    t.integer  "integer1"
    t.integer  "integer2"
    t.integer  "integer3"
    t.integer  "integer4"
    t.integer  "integer5"
    t.integer  "integer6"
    t.integer  "integer7"
    t.integer  "integer8"
    t.integer  "integer9"
    t.integer  "integer10"
    t.integer  "integer11"
    t.integer  "integer12"
    t.integer  "integer13"
    t.integer  "integer14"
    t.integer  "integer15"
    t.integer  "integer16"
    t.integer  "integer17"
    t.integer  "integer18"
    t.integer  "integer19"
    t.integer  "integer20"
    t.float    "float1"
    t.float    "float2"
    t.float    "float3"
    t.float    "float4"
    t.float    "float5"
    t.float    "float6"
    t.float    "float7"
    t.float    "float8"
    t.float    "float9"
    t.float    "float10"
    t.float    "float11"
    t.float    "float12"
    t.float    "float13"
    t.float    "float14"
    t.float    "float15"
    t.float    "float16"
    t.float    "float17"
    t.float    "float18"
    t.float    "float19"
    t.float    "float20"
    t.string   "uri"
  end

  add_index "rows", ["grid_uuid"], :name => "index_rows_on_grid_uuid"
  add_index "rows", ["uri"], :name => "index_rows_on_uri"
  add_index "rows", ["uuid", "begin", "end"], :name => "index_rows_on_uuid_and_begin_and_end"
  add_index "rows", ["uuid"], :name => "index_rows_on_uuid"

  create_table "sessions", :force => true do |t|
    t.string   "session_id",                    :null => false
    t.text     "data",       :limit => 1048576
    t.datetime "created_at"
    t.datetime "updated_at"
  end

  add_index "sessions", ["session_id"], :name => "index_sessions_on_session_id"
  add_index "sessions", ["updated_at"], :name => "index_sessions_on_updated_at"

  create_table "uploads", :force => true do |t|
    t.string   "uuid",             :limit => 36
    t.date     "begin"
    t.date     "end"
    t.integer  "version"
    t.boolean  "enabled"
    t.integer  "lock_version",                   :default => 0
    t.string   "create_user_uuid", :limit => 36
    t.string   "update_user_uuid", :limit => 36
    t.string   "file_name"
    t.datetime "created_at"
    t.datetime "updated_at"
    t.integer  "records"
    t.integer  "inserted"
    t.integer  "updated"
    t.integer  "skipped"
    t.integer  "elapsed"
    t.string   "uri"
  end

  add_index "uploads", ["uri"], :name => "index_uploads_on_uri"

  create_table "users", :force => true do |t|
    t.string   "identifier"
    t.datetime "created_at"
    t.datetime "updated_at"
    t.integer  "lock_version",                         :default => 0
    t.string   "first_name"
    t.string   "last_name"
    t.integer  "version"
    t.date     "begin"
    t.date     "end"
    t.string   "uuid",                   :limit => 36
    t.string   "create_user_uuid",       :limit => 36
    t.string   "update_user_uuid",       :limit => 36
    t.boolean  "enabled"
    t.string   "email",                                :default => "", :null => false
    t.string   "encrypted_password",                   :default => "", :null => false
    t.string   "reset_password_token"
    t.datetime "reset_password_sent_at"
    t.datetime "remember_created_at"
    t.integer  "sign_in_count",                        :default => 0
    t.datetime "current_sign_in_at"
    t.datetime "last_sign_in_at"
    t.string   "current_sign_in_ip"
    t.string   "last_sign_in_ip"
    t.string   "confirmation_token"
    t.datetime "confirmed_at"
    t.datetime "confirmation_sent_at"
    t.string   "unconfirmed_email"
    t.integer  "failed_attempts",                      :default => 0
    t.string   "unlock_token"
    t.datetime "locked_at"
    t.string   "uri"
  end

  add_index "users", ["confirmation_token"], :name => "index_users_on_confirmation_token", :unique => true
  add_index "users", ["email"], :name => "index_users_on_email"
  add_index "users", ["identifier"], :name => "index_users_on_identifier", :unique => true
  add_index "users", ["reset_password_token"], :name => "index_users_on_reset_password_token", :unique => true
  add_index "users", ["unlock_token"], :name => "index_users_on_unlock_token", :unique => true
  add_index "users", ["uri"], :name => "index_users_on_uri"
  add_index "users", ["uuid", "begin", "end"], :name => "index_users_on_uuid_and_begin_and_end"
  add_index "users", ["uuid"], :name => "index_users_on_uuid"

  create_table "workspace_locs", :force => true do |t|
    t.string   "uuid",             :limit => 36
    t.integer  "version"
    t.string   "locale",           :limit => 10
    t.string   "base_locale",      :limit => 10
    t.string   "name"
    t.text     "description"
    t.integer  "lock_version",                   :default => 0
    t.datetime "created_at"
    t.datetime "updated_at"
    t.string   "create_user_uuid", :limit => 36
    t.string   "update_user_uuid", :limit => 36
  end

  add_index "workspace_locs", ["uuid", "version", "locale"], :name => "index_workspace_locs_on_uuid_and_version_and_locale"

  create_table "workspace_sharings", :force => true do |t|
    t.string   "uuid",             :limit => 36
    t.date     "begin"
    t.date     "end"
    t.integer  "version"
    t.boolean  "enabled"
    t.string   "workspace_uuid",   :limit => 36
    t.string   "user_uuid",        :limit => 36
    t.string   "role_uuid",        :limit => 36
    t.integer  "lock_version",                   :default => 0
    t.string   "create_user_uuid", :limit => 36
    t.string   "update_user_uuid", :limit => 36
    t.datetime "created_at",                                    :null => false
    t.datetime "updated_at",                                    :null => false
    t.string   "uri"
  end

  add_index "workspace_sharings", ["uri"], :name => "index_workspace_sharings_on_uri"
  add_index "workspace_sharings", ["uuid", "begin", "end"], :name => "index_workspace_sharings_on_uuid_and_begin_and_end"
  add_index "workspace_sharings", ["workspace_uuid", "user_uuid"], :name => "index_workspace_sharings_on_workspace_uuid_and_user_uuid"

  create_table "workspaces", :force => true do |t|
    t.date     "begin"
    t.date     "end"
    t.integer  "version"
    t.datetime "created_at"
    t.datetime "updated_at"
    t.integer  "lock_version",                    :default => 0
    t.string   "uuid",              :limit => 36
    t.string   "create_user_uuid",  :limit => 36
    t.string   "update_user_uuid",  :limit => 36
    t.boolean  "enabled"
    t.boolean  "public"
    t.string   "default_role_uuid", :limit => 36
    t.string   "uri"
  end

  add_index "workspaces", ["public"], :name => "index_workspaces_on_public"
  add_index "workspaces", ["uri"], :name => "index_workspaces_on_uri"
  add_index "workspaces", ["uuid", "begin", "end"], :name => "index_workspaces_on_uuid_and_begin_and_end"
  add_index "workspaces", ["uuid"], :name => "index_workspaces_on_uuid"

end
