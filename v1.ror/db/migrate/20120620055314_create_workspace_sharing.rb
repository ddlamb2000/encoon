class CreateWorkspaceSharing < ActiveRecord::Migration
  def self.up
    create_table :workspace_sharings do |t|
      t.string   :uuid,                :limit => 36
      t.date     :begin
      t.date     :end
      t.integer  :version
      t.boolean  :enabled
      t.string   :workspace_uuid,    :limit => 36
      t.string   :user_uuid,         :limit => 36
      t.string   :role_uuid,         :limit => 36
      t.integer  :lock_version,        :default => 0
      t.string   :create_user_uuid,    :limit => 36
      t.string   :update_user_uuid,    :limit => 36
      t.timestamps
    end

    add_index :workspace_sharings, [:uuid, :begin, :end]
    add_index :workspace_sharings, [:workspace_uuid, :user_uuid]
  end

  def self.down
    drop_table :workspace_sharings
  end
end