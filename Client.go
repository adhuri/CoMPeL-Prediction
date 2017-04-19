package main

import (
	"encoding/gob"
	"net"

	"github.com/Sirupsen/logrus"
	"github.com/adhuri/Compel-Migration/protocol"
)

func sendDataToMigration(dataToSendToMigration *protocol.PredictionData, log *logrus.Logger) error {
	addr := "127.0.0.1" + ":" + "5051"
	log.Info("Location of Migration Server : ", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Errorln("Migration Server Not Alive")
		return err
	}
	// If connection successful send a connect message
	encoder := gob.NewEncoder(conn)
	err = encoder.Encode(dataToSendToMigration)
	if err != nil {
		log.Errorln("Failure While Sending Data To Server ")
		return err
	}
	log.Infoln("Prediction Data Message Successfully Sent")

	// read ack from the server
	serverReply := protocol.PredictionDataResponse{}
	decoder := gob.NewDecoder(conn)
	err = decoder.Decode(&serverReply)
	if err != nil {
		// If error occurs while reading ACK from server then return
		log.Errorln("ACK for Prediction Data Message Failed " + err.Error())
		return err
	} // Print the ACK received from the server
	log.Infoln("Prediction Data Message ACK Received")
	return nil
}
