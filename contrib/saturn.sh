#!/usr/bin/env fish

set -l options (fish_opt -s h -l help)
set options $options (fish_opt -s m -l mode --required-val)
set options $options (fish_opt -s n -l num --required-val)
argparse $options -- $argv

if set -q _flag_help
    saturn -help
    return 0
end

if set -q _flag_mode and set -q _flag_num
    if set x (saturn -search $argv | fzf)
        saturn -fetch $x $_flag_mode $_flag_num
    end
else
    if set x (saturn -search $argv | fzf)
        set num_ep (saturn -fetch $x | fzf -m | grep -oP 'ID:\s*\K\d+' | tr ' ' '-')
		set sanitized_num_ep (echo $num_ep | tr ' ' '-')
	    saturn -fetch $x $_flag_mode $sanitized_num_ep
        return 0
    end
end