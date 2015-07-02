#!/usr/bin/env ruby

require 'tmpdir'

dir = Dir.mktmpdir
Dir.chdir(dir)

`got init`

File.write('file1.txt', 'First file')
File.write('file2.txt', 'Second file')

result = `got status`

expected = %q(Comparing working directory to empty tree.
A file1.txt
A file2.txt
)

if result != expected
  exit 1
end
