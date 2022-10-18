class CreateRows < ActiveRecord::Migration
  def self.up
    create_table :rows do |t|
      t.date :begin
      t.date :end
      t.integer :version
      t.integer :lock_version
      t.references :grid
      t.integer :create_user_id
      t.integer :update_user_id
      t.string :name
      t.text :description
      t.string :value1
      t.string :value2
      t.string :value3
      t.string :value4
      t.string :value5
      t.string :value6
      t.string :value7
      t.string :value8
      t.string :value9
      t.string :value10
      t.string :value11
      t.string :value12
      t.string :value13
      t.string :value14
      t.string :value15
      t.string :value16
      t.string :value17
      t.string :value18
      t.string :value19
      t.string :value20

      t.timestamps
    end
  end

  def self.down
    drop_table :rows
  end
end
