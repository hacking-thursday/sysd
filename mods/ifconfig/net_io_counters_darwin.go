package ifconfig

func net_io_counters() (iface_counter_map map[string] Counters) {
    iface_counter_map = make(map[string]Counters)

    return iface_counter_map
}
