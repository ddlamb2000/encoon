class RemoveLengthDecimalsColumn < ActiveRecord::Migration
  def self.up
    remove_column :columns, :length
    remove_column :columns, :decimals
  end

  def self.down
    add_column :columns, :length, :integer
    add_column :columns, :decimals, :integer
  end
end
