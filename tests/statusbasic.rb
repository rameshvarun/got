#!/usr/bin/env ruby

require 'tmpdir'

expected = IO.read(File.join(File.dirname(__FILE__), "statusbasic.out"))

dir = Dir.mktmpdir
Dir.chdir(dir)

`got init`

File.write('file1.txt', 'First file')
File.write('file2.txt', 'Second file')

`got commit -m "Test Message" -a "Test Author"`

File.write('file3.txt', 'Third file')
File.write('file2.txt', 'Modified second file')
File.delete('file1.txt')

result = `got status`

if result != expected
  exit 1
end
