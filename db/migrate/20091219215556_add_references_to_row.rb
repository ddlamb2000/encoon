class AddReferencesToRow < ActiveRecord::Migration
  def self.up
    add_column :rows, :row_id1, :integer
    add_column :rows, :row_id2, :integer
    add_column :rows, :row_id3, :integer
    add_column :rows, :row_id4, :integer
    add_column :rows, :row_id5, :integer
    add_column :rows, :row_id6, :integer
    add_column :rows, :row_id7, :integer
    add_column :rows, :row_id8, :integer
    add_column :rows, :row_id9, :integer
    add_column :rows, :row_id10, :integer
    add_column :rows, :row_id11, :integer
    add_column :rows, :row_id12, :integer
    add_column :rows, :row_id13, :integer
    add_column :rows, :row_id14, :integer
    add_column :rows, :row_id15, :integer
    add_column :rows, :row_id16, :integer
    add_column :rows, :row_id17, :integer
    add_column :rows, :row_id18, :integer
    add_column :rows, :row_id19, :integer
    add_column :rows, :row_id20, :integer
  end

  def self.down
    remove_column :rows, :row_id1
    remove_column :rows, :row_id2
    remove_column :rows, :row_id3
    remove_column :rows, :row_id4
    remove_column :rows, :row_id5
    remove_column :rows, :row_id6
    remove_column :rows, :row_id7
    remove_column :rows, :row_id8
    remove_column :rows, :row_id9
    remove_column :rows, :row_id10
    remove_column :rows, :row_id11
    remove_column :rows, :row_id12
    remove_column :rows, :row_id13
    remove_column :rows, :row_id14
    remove_column :rows, :row_id15
    remove_column :rows, :row_id16
    remove_column :rows, :row_id17
    remove_column :rows, :row_id18
    remove_column :rows, :row_id19
    remove_column :rows, :row_id20
  end
end
