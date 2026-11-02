package monitor

// probeScript prints everything a Reading needs in one round trip, as
// sections introduced by "@@name" lines. It's POSIX sh and reads /proc
// and /sys, so it works on any Linux, BusyBox included, without
// installing anything on the server.
const probeScript = `LC_ALL=C; export LC_ALL
echo @@host; hostname 2>/dev/null || cat /proc/sys/kernel/hostname
echo @@os; ( . /etc/os-release 2>/dev/null && echo "$PRETTY_NAME" ) || echo; uname -srm
echo @@uptime; cat /proc/uptime
echo @@loadavg; cat /proc/loadavg
echo @@cpuinfo; grep -m1 -i -E '^(model name|hardware|cpu model)' /proc/cpuinfo
echo @@stat; grep '^cpu' /proc/stat
echo @@meminfo; cat /proc/meminfo
echo @@df; df -PTk 2>/dev/null
echo @@diskstats; cat /proc/diskstats 2>/dev/null
echo @@netdev; cat /proc/net/dev
echo @@thermal; for z in /sys/class/thermal/thermal_zone*; do [ -r "$z/temp" ] && echo "$(cat "$z/type" 2>/dev/null)|$(cat "$z/temp" 2>/dev/null)"; done
echo @@hwmon; for h in /sys/class/hwmon/hwmon*; do n=$(cat "$h/name" 2>/dev/null); for t in "$h"/temp*_input; do [ -r "$t" ] || continue; b=${t%_input}; echo "$n|$(cat "${b}_label" 2>/dev/null)|$(cat "$t" 2>/dev/null)|$(cat "${b}_max" 2>/dev/null)|$(cat "${b}_crit" 2>/dev/null)"; done; done
echo @@procs; ls -d /proc/[0-9]* 2>/dev/null | wc -l
echo @@end
`
