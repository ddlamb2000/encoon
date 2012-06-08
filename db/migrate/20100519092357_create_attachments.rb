class CreateAttachments < ActiveRecord::Migration
  def self.up
    create_table :row_attachments do |t|
      t.string   :uuid, :limit => 36
      t.string   :file_name
      t.binary   :document, :limit => 1048576
      t.string   :content_type
      t.integer  :lock_version, :default => 0
    end

    add_index :row_attachments, [:uuid]
  end

  def self.down
    drop_table :row_attachments
  end
end
