class AddCountToUpload < ActiveRecord::Migration
  def self.up
    add_column :uploads, :records, :integer
    add_column :uploads, :inserted, :integer
    add_column :uploads, :updated, :integer
  end

  def self.down
    remove_column :uploads, :records
    remove_column :uploads, :inserted
    remove_column :uploads, :updated
  end
end
