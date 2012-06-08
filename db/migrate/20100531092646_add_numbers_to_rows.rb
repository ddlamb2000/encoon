class AddNumbersToRows < ActiveRecord::Migration
  def self.up
    add_column :rows, :integer1, :integer
    add_column :rows, :integer2, :integer
    add_column :rows, :integer3, :integer
    add_column :rows, :integer4, :integer
    add_column :rows, :integer5, :integer
    add_column :rows, :integer6, :integer
    add_column :rows, :integer7, :integer
    add_column :rows, :integer8, :integer
    add_column :rows, :integer9, :integer
    add_column :rows, :integer10, :integer
    add_column :rows, :integer11, :integer
    add_column :rows, :integer12, :integer
    add_column :rows, :integer13, :integer
    add_column :rows, :integer14, :integer
    add_column :rows, :integer15, :integer
    add_column :rows, :integer16, :integer
    add_column :rows, :integer17, :integer
    add_column :rows, :integer18, :integer
    add_column :rows, :integer19, :integer
    add_column :rows, :integer20, :integer

    add_column :rows, :float1, :float
    add_column :rows, :float2, :float
    add_column :rows, :float3, :float
    add_column :rows, :float4, :float
    add_column :rows, :float5, :float
    add_column :rows, :float6, :float
    add_column :rows, :float7, :float
    add_column :rows, :float8, :float
    add_column :rows, :float9, :float
    add_column :rows, :float10, :float
    add_column :rows, :float11, :float
    add_column :rows, :float12, :float
    add_column :rows, :float13, :float
    add_column :rows, :float14, :float
    add_column :rows, :float15, :float
    add_column :rows, :float16, :float
    add_column :rows, :float17, :float
    add_column :rows, :float18, :float
    add_column :rows, :float19, :float
    add_column :rows, :float20, :float
  end

  def self.down
    remove_column :rows, :integer1
    remove_column :rows, :integer2
    remove_column :rows, :integer3
    remove_column :rows, :integer4
    remove_column :rows, :integer5
    remove_column :rows, :integer6
    remove_column :rows, :integer7
    remove_column :rows, :integer8
    remove_column :rows, :integer9
    remove_column :rows, :integer10
    remove_column :rows, :integer11
    remove_column :rows, :integer12
    remove_column :rows, :integer13
    remove_column :rows, :integer14
    remove_column :rows, :integer15
    remove_column :rows, :integer16
    remove_column :rows, :integer17
    remove_column :rows, :integer18
    remove_column :rows, :integer19
    remove_column :rows, :integer20

    remove_column :rows, :float1
    remove_column :rows, :float2
    remove_column :rows, :float3
    remove_column :rows, :float4
    remove_column :rows, :float5
    remove_column :rows, :float6
    remove_column :rows, :float7
    remove_column :rows, :float8
    remove_column :rows, :float9
    remove_column :rows, :float10
    remove_column :rows, :float11
    remove_column :rows, :float12
    remove_column :rows, :float13
    remove_column :rows, :float14
    remove_column :rows, :float15
    remove_column :rows, :float16
    remove_column :rows, :float17
    remove_column :rows, :float18
    remove_column :rows, :float19
    remove_column :rows, :float20
  end
end
