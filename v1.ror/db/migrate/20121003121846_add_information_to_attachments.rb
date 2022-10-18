class AddInformationToAttachments < ActiveRecord::Migration
  def self.up
    add_column :attachments, :original_file_name, :string
    add_column :attachments, :create_user_uuid, :string, :limit => 36
  end

  def self.down
    remove_column :attachments, :original_file_name
    remove_column :attachments, :create_user_uuid
  end
end