class UpdateFkUuidData < ActiveRecord::Migration
  def self.up
    Grid.find(:all).each do |item|
      puts "Update grid " + item.id.to_s
      if item.workspace_id.present?
        item.workspace_uuid = Workspace.find(item.workspace_id).uuid
        item.save! if item.valid?
      end
      puts "...create_user_id=" + item.create_user_id.to_s
      item.create_user_uuid = User.find(item.create_user_id).uuid
      puts "......create_user_uuid=" + item.create_user_uuid.to_s
      item.save! if item.valid?
      puts "...update_user_id=" + item.update_user_id.to_s
      item.update_user_uuid = User.find(item.update_user_id).uuid
      puts "......update_user_uuid=" + item.update_user_uuid.to_s
      item.save! if item.valid?
    end

    Column.find(:all).each do |item|
      puts "Update column " + item.id.to_s
      if item.grid_id.present?
        puts "...grid_id=" + item.grid_id.to_s
        item.grid_uuid = Grid.find(item.grid_id).uuid
        item.save! if item.valid?
      end
      if item.grid_reference_id.present?
        puts "...grid_reference_id=" + item.grid_reference_id.to_s
        item.grid_reference_uuid = Grid.find(item.grid_reference_id).uuid
        item.save! if item.valid?
      end
      puts "...create_user_id=" + item.create_user_id.to_s
      item.create_user_uuid = User.find(item.create_user_id).uuid
      puts "......create_user_uuid=" + item.create_user_uuid.to_s
      item.save! if item.valid?
      puts "...update_user_id=" + item.update_user_id.to_s
      item.update_user_uuid = User.find(item.update_user_id).uuid
      puts "......update_user_uuid=" + item.update_user_uuid.to_s
      item.save! if item.valid?
    end

    Row.find(:all).each do |item|
      puts "Update row " + item.id.to_s
      if item.grid_id.present?
        puts "...grid_id=" + item.grid_id.to_s
        item.grid_uuid = Grid.find(item.grid_id).uuid
        if item.valid?
          item.save!
        else
          puts "Record isn't valid!"
        end
      end
      if item.row_id1.present?
        puts "...row_id1=" + item.row_id1.to_s
        item.row_uuid1 = Row.find(item.row_id1).uuid
        if item.valid?
          item.save!
        else
          puts "Record isn't valid!"
        end
      end
      if item.row_id2.present?
        puts "...row_id2=" + item.row_id2.to_s
        item.row_uuid2 = Row.find(item.row_id2).uuid
        if item.valid?
          item.save!
        else
          puts "Record isn't valid!"
        end
      end
      if item.row_id3.present?
        puts "...row_id3=" + item.row_id3.to_s
        item.row_uuid3 = Row.find(item.row_id3).uuid
        if item.valid?
          item.save!
        else
          puts "Record isn't valid!"
        end
      end
      if item.row_id4.present?
        puts "...row_id4=" + item.row_id4.to_s
        item.row_uuid4 = Row.find(item.row_id4).uuid
        if item.valid?
          item.save!
        else
          puts "Record isn't valid!"
        end
      end
      if item.row_id5.present?
        puts "...row_id5=" + item.row_id5.to_s
        item.row_uuid5 = Row.find(item.row_id5).uuid
        if item.valid?
          item.save!
        else
          puts "Record isn't valid!"
        end
      end
      if item.row_id6.present?
        puts "...row_id6=" + item.row_id6.to_s
        item.row_uuid6 = Row.find(item.row_id6).uuid
        if item.valid?
          item.save!
        else
          puts "Record isn't valid!"
        end
      end
      if item.row_id7.present?
        puts "...row_id7=" + item.row_id7.to_s
        item.row_uuid7 = Row.find(item.row_id7).uuid
        if item.valid?
          item.save!
        else
          puts "Record isn't valid!"
        end
      end
      if item.row_id8.present?
        puts "...row_id8=" + item.row_id8.to_s
        item.row_uuid8 = Row.find(item.row_id8).uuid
        if item.valid?
          item.save!
        else
          puts "Record isn't valid!"
        end
      end
      puts "...create_user_id=" + item.create_user_id.to_s
      item.create_user_uuid = User.find(item.create_user_id).uuid
      puts "......create_user_uuid=" + item.create_user_uuid.to_s
      item.save! if item.valid?
      puts "...update_user_id=" + item.update_user_id.to_s
      item.update_user_uuid = User.find(item.update_user_id).uuid
      puts "......update_user_uuid=" + item.update_user_uuid.to_s
      item.save! if item.valid?
    end

    User.find(:all).each do |item|
      puts "Update user " + item.id.to_s
      puts "...create_user_id=" + item.create_user_id.to_s
      item.create_user_uuid = User.find(item.create_user_id).uuid
      puts "......create_user_uuid=" + item.create_user_uuid.to_s
      item.save! if item.valid?
      puts "...update_user_id=" + item.update_user_id.to_s
      item.update_user_uuid = User.find(item.update_user_id).uuid
      puts "......update_user_uuid=" + item.update_user_uuid.to_s
      item.save! if item.valid?
    end

    Workspace.find(:all).each do |item|
      puts "Update workspace " + item.id.to_s
      puts "...create_user_id=" + item.create_user_id.to_s
      item.create_user_uuid = User.find(item.create_user_id).uuid
      puts "......create_user_uuid=" + item.create_user_uuid.to_s
      item.save! if item.valid?
      puts "...update_user_id=" + item.update_user_id.to_s
      item.update_user_uuid = User.find(item.update_user_id).uuid
      puts "......update_user_uuid=" + item.update_user_uuid.to_s
      item.save! if item.valid?
    end
  end

  def self.down
  end
end
