class CreateUploads < ActiveRecord::Migration
  def self.up
    create_table :uploads do |t|
      t.string   :uuid,                :limit => 36
      t.date     :begin
      t.date     :end
      t.integer  :version
      t.boolean  :enabled
      t.integer  :lock_version,        :default => 0
      t.string   :create_user_uuid,    :limit => 36
      t.string   :update_user_uuid,    :limit => 36
      t.string   :file_name
      t.timestamps
    end
  end

  def self.down
    drop_table :uploads
  end
end
