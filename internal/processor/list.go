package processor

import (
	"fmt"

	"github.com/donkeysharp/donkeyvpn/internal/service"
	"github.com/donkeysharp/donkeyvpn/internal/telegram"
	"github.com/labstack/gommon/log"
)

func NewListProcessor(client *telegram.Client, vpnService *service.VPNService, peerService *service.PeerService) ListProcessor {
	return ListProcessor{
		ProcessorShared: ProcessorShared{
			Client: client,
		},
		vpnService:  vpnService,
		peerService: peerService,
	}
}

type ListProcessor struct {
	ProcessorShared
	vpnService  *service.VPNService
	peerService *service.PeerService
}

func (p ListProcessor) sendMessage(msg string, update *telegram.Update) {
	err := p.Client.SendMessage(msg, update.Message.Chat)
	if err != nil {
		log.Errorf("Error sending message to Telegram. msg=%s", msg)
	}
}

func (p ListProcessor) ListVPNs(update *telegram.Update) error {
	log.Info("Listing vpn instances for telegram")
	instances, err := p.vpnService.ListArray()
	if err != nil {
		p.sendMessage("Error while retrieving VPN instances. Try again please.", update)
	}
	message := "List of instances:\n-----\n"
	for _, item := range instances {
		log.Infof("Proccessing instance: %s %s %s %s", item.Id, item.Hostname, item.Port, item.Status)
		message += fmt.Sprintf("Id: %s\n", item.Id)
		message += fmt.Sprintf("Hostname: %s\n", item.Hostname)
		message += fmt.Sprintf("Port: %s\n", item.Port)
		message += fmt.Sprintf("Status: %s\n", item.Status)
		message += "-----\n"
	}

	p.sendMessage(message, update)
	return nil
}

func (p ListProcessor) ListPeers(update *telegram.Update) error {
	log.Infof("Listing peers for telegram")
	peers, err := p.peerService.List()
	if err != nil {
		p.sendMessage("Error while retrieving peers. Try again please.", update)
		return err
	}

	message := "List of peers:\n-----\n"
	for _, item := range peers {
		log.Infof("Processing peer: %s", item.IPAddress)
		message += fmt.Sprintf("IP Address: %s\n", item.IPAddress)
		message += fmt.Sprintf("Public Key: %s\n", item.PublicKey)
		message += fmt.Sprintf("Username: %s\n", item.Username)
		message += "-----\n"
	}

	p.sendMessage(message, update)
	return nil
}

func (p ListProcessor) Process(args []string, update *telegram.Update) error {
	log.Infof("Processing '/list' command with args %v for chat %d", args, update.Message.Chat.ChatId)

	usage := getUsage()
	if len(args) >= 1 && args[0] == "vpn" {
		return p.ListVPNs(update)
	}

	if len(args) >= 1 && args[0] == "peers" {
		return p.ListPeers(update)
	}

	err := p.Client.SendMessage(usage, update.Message.Chat)
	if err != nil {
		log.Errorf("Error sending message to Telegram. msg=%s", usage)
	}
	return nil
}
