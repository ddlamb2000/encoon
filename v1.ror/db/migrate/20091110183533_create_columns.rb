class CreateColumns < ActiveRecord::Migration
  def self.up
    create_table :columns do |t|
      t.date :begin
      t.date :end
      t.string :name
      t.text :description
      t.integer :version
      t.integer :revision
      t.references :grid
      t.integer :number
      t.integer :display
      t.string :kind
      t.boolean :translated

      t.timestamps
    end
  end

  def self.down
    drop_table :columns
  end
end
