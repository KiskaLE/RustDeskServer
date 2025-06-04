package utils

import (
	"log"
	"os/exec"
)

func RemoveFirewallRule(name string) error {
	cmd := exec.Command("netsh", "advfirewall", "firewall", "delete", "rule", "name="+name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Chyba při mazání firewall pravidla: %s\n", output)
		return err
	}
	return nil
}

func ConfigureFirewall(name string, port string) error {
	// Povolení portu ve Windows Firewall (Inbound)
	cmd := exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
		"name="+name,
		"dir=in",
		"action=allow",
		"protocol=TCP",
		"localport="+port,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Nepodařilo se přidat firewall pravidlo: %s\n", output)
		return err
	}
	return nil
}
