class CreateRolesAndAuthorizations < ActiveRecord::Migration
  def self.up
    create_table :roles do |t|
      t.string   :uuid,                :limit => 36
      t.date     :begin
      t.date     :end
      t.integer  :version
      t.boolean  :enabled
      t.integer  :lock_version,        :default => 0
      t.string   :create_user_uuid,    :limit => 36
      t.string   :update_user_uuid,    :limit => 36
      t.timestamps
    end

    add_index :roles, [:uuid, :begin, :end]
    add_index :roles, [:uuid]

    create_table :role_locs do |t|
      t.string   :uuid, :limit => 36
      t.integer  :version
      t.string   :locale, :limit => 10
      t.string   :base_locale, :limit => 10
      t.string   :name
      t.text     :description
      t.integer  :lock_version, :default => 0
    end

    add_index :role_locs, [:uuid, :version, :locale]

    create_table :user_roles do |t|
      t.string   :uuid,                :limit => 36
      t.date     :begin
      t.date     :end
      t.integer  :version
      t.boolean  :enabled
      t.string   :user_uuid,    :limit => 36
      t.string   :role_uuid,    :limit => 36
      t.integer  :lock_version,        :default => 0
      t.string   :create_user_uuid,    :limit => 36
      t.string   :update_user_uuid,    :limit => 36
      t.timestamps
    end

    add_index :user_roles, [:uuid, :begin, :end, :user_uuid, :role_uuid]
    add_index :user_roles, [:uuid, :begin, :end]
    add_index :user_roles, [:uuid]

    create_table :grid_authorizations do |t|
      t.string   :uuid,                :limit => 36
      t.date     :begin
      t.date     :end
      t.integer  :version
      t.boolean  :enabled
      t.string   :grid_uuid,    :limit => 36
      t.string   :role_uuid,    :limit => 36
      t.boolean  :grid_select
      t.boolean  :grid_update
      t.boolean  :grid_delete
      t.boolean  :data_select
      t.boolean  :data_create
      t.boolean  :data_update
      t.boolean  :data_delete
      t.integer  :lock_version,        :default => 0
      t.string   :create_user_uuid,    :limit => 36
      t.string   :update_user_uuid,    :limit => 36
      t.timestamps
    end

    add_index :grid_authorizations, [:uuid, :begin, :end, :grid_uuid, :role_uuid]
    add_index :grid_authorizations, [:uuid, :begin, :end]
    add_index :grid_authorizations, [:uuid]
  end

  def self.down
    drop_table :roles
    drop_table :role_locs
    drop_table :user_roles
    drop_table :grid_authorizations
  end
end
