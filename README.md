# *HYPE* Script
- A Hype scripting language

### To-Do
- Async scripting, keeping it simple
- Easy importing
- Basically if Python and Go had a child, but Bash watched from the closet crack

### What problems does this solve?
- Async not easily in bash
- Bash not 'beautiful'
- 

### Endgoal
- Usable replacement for bash on my system
- Re-write VPNBeans


For instance:
"""
down() {
    stop_openvpn
    sleep 1s
    del_routes
    stop_stunnel
	rm_dns
}
"""

// I would like to stop_openvpn, del_routes & stop_stunnel all at once. We need as many to stop as possible so I don't care if one doesnt stop, they all NEED to

"""
fun down() {
    x = par { // What does this return? List of returned values?
        stop_openvpn, // Notice how I can call functions that dont have any params without ()
        del_routes,
        stop_stunnel,
        rm_dns,
    }
    // Goes to Stderr if failure
}
"""