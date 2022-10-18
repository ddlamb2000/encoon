class AddPaperClipToAttachments < ActiveRecord::Migration
  def self.up
    drop_table :row_attachments

    create_table :attachments do |t|
      t.string   :uuid, :limit => 36
      t.string   :document_file_name
      t.string   :document_content_type
      t.integer  :document_file_size
      t.datetime :document_updated_at
      t.integer  :lock_version, :default => 0
    end

    add_index :attachments, [:uuid]
  end

  def self.down
  end
end