class RemoveTranslatedColumn < ActiveRecord::Migration
  def self.up
    remove_column :columns, :translated
  end

  def self.down
    add_column :columns, :translated, :boolean
  end
end
