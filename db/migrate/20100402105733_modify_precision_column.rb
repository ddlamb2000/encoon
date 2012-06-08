class ModifyPrecisionColumn < ActiveRecord::Migration
  def self.up
    remove_column :columns, :precision
    add_column :columns, :decimals, :integer
  end

  def self.down
    add_column :columns, :precision, :integer
    remove_column :decimals, :precision
  end
end
