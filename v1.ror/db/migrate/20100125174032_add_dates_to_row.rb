class AddDatesToRow < ActiveRecord::Migration
  def self.up
    add_column :rows, :date1, :date
    add_column :rows, :date2, :date
    add_column :rows, :date3, :date
    add_column :rows, :date4, :date
    add_column :rows, :date5, :date
    add_column :rows, :date6, :date
    add_column :rows, :date7, :date
    add_column :rows, :date8, :date
    add_column :rows, :date9, :date
    add_column :rows, :date10, :date
    add_column :rows, :date11, :date
    add_column :rows, :date12, :date
    add_column :rows, :date13, :date
    add_column :rows, :date14, :date
    add_column :rows, :date15, :date
    add_column :rows, :date16, :date
    add_column :rows, :date17, :date
    add_column :rows, :date18, :date
    add_column :rows, :date19, :date
    add_column :rows, :date20, :date
  end

  def self.down
    remove_column :rows, :date1
    remove_column :rows, :date2
    remove_column :rows, :date3
    remove_column :rows, :date4
    remove_column :rows, :date5
    remove_column :rows, :date6
    remove_column :rows, :date7
    remove_column :rows, :date8
    remove_column :rows, :date9
    remove_column :rows, :date10
    remove_column :rows, :date11
    remove_column :rows, :date12
    remove_column :rows, :date13
    remove_column :rows, :date14
    remove_column :rows, :date15
    remove_column :rows, :date16
    remove_column :rows, :date17
    remove_column :rows, :date18
    remove_column :rows, :date19
    remove_column :rows, :date20
  end
end
