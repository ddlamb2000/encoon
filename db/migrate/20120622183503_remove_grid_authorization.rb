class RemoveGridAuthorization < ActiveRecord::Migration
  def up
    drop_table :user_roles
    drop_table :grid_authorizations
  end

  def down
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
end
