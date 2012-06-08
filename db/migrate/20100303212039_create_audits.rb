class CreateAudits < ActiveRecord::Migration
  def self.up
    create_table :audits do |t|
      t.string   :uuid,                :limit => 36
      t.date     :begin
      t.date     :end
      t.integer  :version
      t.boolean  :enabled
      t.integer  :lock_version,        :default => 0
      t.string   :create_user_uuid,    :limit => 36
      t.string   :update_user_uuid,    :limit => 36
      t.string   :grid_uuid,           :limit => 36
      t.string   :row_uuid,            :limit => 36
      t.string   :kind,                :limit => 36
      t.timestamps
    end

    add_index :audits, [:grid_uuid, :row_uuid, :version]
  end

  def self.down
    drop_table :audits
  end
end
