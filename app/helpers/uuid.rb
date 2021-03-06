#
# = uuid.rb - UUID generator
#
# Author:: Assaf Arkin  assaf@labnotes.org
#          Eric Hodel drbrain@segment7.net
# Copyright:: Copyright (c) 2005-2010 Assaf Arkin, Eric Hodel
# License:: MIT and/or Creative Commons Attribution-ShareAlike

require 'thread'
require 'macaddr'

##
# = Generating UUIDs
#
# Call #generate to generate a new UUID. The method returns a string in one of
# three formats. The default format is 36 characters long, and contains the 32
# hexadecimal octets and hyphens separating the various value parts. The
# <tt>:compact</tt> format omits the hyphens, while the <tt>:urn</tt> format
# adds the <tt>:urn:uuid</tt> prefix.
#
# For example:
#
#   uuid = UUID.new
#   
#   10.times do
#     p uuid.generate
#   end
#
# = UUIDs in Brief
#
# UUID (universally unique identifier) are guaranteed to be unique across time
# and space.
#
# A UUID is 128 bit long, and consists of a 60-bit time value, a 16-bit
# sequence number and a 48-bit node identifier.
#
# The time value is taken from the system clock, and is monotonically
# incrementing.  However, since it is possible to set the system clock
# backward, a sequence number is added.  The sequence number is incremented
# each time the UUID generator is started.  The combination guarantees that
# identifiers created on the same machine are unique with a high degree of
# probability.
#
# Note that due to the structure of the UUID and the use of sequence number,
# there is no guarantee that UUID values themselves are monotonically
# incrementing.  The UUID value cannot itself be used to sort based on order
# of creation.
#
# To guarantee that UUIDs are unique across all machines in the network,
# the IEEE 802 MAC address of the machine's network interface card is used as
# the node identifier.
#
# For more information see {RFC 4122}[http://www.ietf.org/rfc/rfc4122.txt].

class UUID

  ##
  # Clock multiplier. Converts Time (resolution: seconds) to UUID clock
  # (resolution: 10ns)
  CLOCK_MULTIPLIER = 10000000

  ##
  # Clock gap is the number of ticks (resolution: 10ns) between two Ruby Time
  # ticks.
  CLOCK_GAPS = 100000

  ##
  # Version number stamped into the UUID to identify it as time-based.
  VERSION_CLOCK = 0x0100

  ##
  # Formats supported by the UUID generator.
  #
  # <tt>:default</tt>:: Produces 36 characters, including hyphens separating
  #                     the UUID value parts
  # <tt>:compact</tt>:: Produces a 32 digits (hexadecimal) value with no
  #                     hyphens
  # <tt>:urn</tt>:: Adds the prefix <tt>urn:uuid:</tt> to the default format
  FORMATS = {
    :compact => '%08x%04x%04x%04x%012x',
    :default => '%08x-%04x-%04x-%04x-%012x',
    :urn     => 'urn:uuid:%08x-%04x-%04x-%04x-%012x',
  }

  @uuid = nil

  ##
  # Generates a new UUID string using +format+.  See FORMATS for a list of
  # supported formats.

  def self.generate(format = :default)
    @uuid ||= new
    @uuid.generate format
  end

  ##
  # Returns true if +uuid+ is in compact, default or urn formats.  Does not
  # validate the layout (RFC 4122 section 4) of the UUID.
  def self.validate(uuid)
    return true if uuid =~ /\A[\da-f]{32}\z/i
    return true if
      uuid =~ /\A(urn:uuid:)?[\da-f]{8}-([\da-f]{4}-){3}[\da-f]{12}\z/i
  end

  ##
  # Create a new UUID generator.  You really only need to do this once.
  def initialize
    @drift = 0
    @last_clock = (Time.now.to_f * CLOCK_MULTIPLIER).to_i
    @mutex = Mutex.new

    @mac = Mac.addr.gsub(/:|-/, '').hex & 0x7FFFFFFFFFFF
    fail "Cannot determine MAC address from any available interface, tried with #{Mac.addr}" if @mac == 0
    @sequence = rand 0x10000
  end

  ##
  # Generates a new UUID string using +format+.  See FORMATS for a list of
  # supported formats.
  def generate(format = :default)
    template = FORMATS[format]

    raise ArgumentError, "invalid UUID format #{format.inspect}" unless template

    # The clock must be monotonically increasing. The clock resolution is at
    # best 100 ns (UUID spec), but practically may be lower (on my setup,
    # around 1ms). If this method is called too fast, we don't have a
    # monotonically increasing clock, so the solution is to just wait.
    #
    # It is possible for the clock to be adjusted backwards, in which case we
    # would end up blocking for a long time. When backward clock is detected,
    # we prevent duplicates by asking for a new sequence number and continue
    # with the new clock.

    clock = @mutex.synchronize do
      clock = (Time.new.to_f * CLOCK_MULTIPLIER).to_i & 0xFFFFFFFFFFFFFFF0

      if clock > @last_clock then
        @drift = 0
        @last_clock = clock
      elsif clock == @last_clock then
        drift = @drift += 1

        if drift < 10000 then
          @last_clock += 1
        else
          Thread.pass
          nil
        end
      else
        next_sequence
        @last_clock = clock
      end
    end until clock

    template % [
        clock        & 0xFFFFFFFF,
       (clock >> 32) & 0xFFFF,
      ((clock >> 48) & 0xFFFF | VERSION_CLOCK),
      @sequence      & 0xFFFF,
      @mac           & 0xFFFFFFFFFFFF
    ]
  end

  # Updates with a new sequence number.
  def next_sequence
    @sequence += 1
    @last_clock = (Time.now.to_f * CLOCK_MULTIPLIER).to_i
    @drift = 0
  end

  def inspect
    mac = ("%012x" % @mac).scan(/[0-9a-f]{2}/).join(':')
    "MAC: #{mac}  Sequence: #{@sequence}"
  end
end
