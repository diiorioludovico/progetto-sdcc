import grpc
import advancedSensorSimulator
import json
import datetime
import time
import edge_pb2
import edge_pb2_grpc
import traceback
import os
import logging

thresholds = {"temperature": 0.5,
              "humidity": 2,
              "brightness": 20,
              "air_quality": 5}

logger = logging.getLogger("EdgeDevice")

class ParkSensor:

    def __init__(self, server, serial_number):
        self.server = server
        self.serial_number = serial_number
        logger.info("Sensor initialization")

        #initialization comunication with server
        result = self.init()

        self.deviceID = result.deviceID
        self.parkID = result.parkID
        self.interval = result.interval

        logger.info("Edge device ID: " + str(self.deviceID))
        logger.info("Associated park ID: " + str(self.parkID))

        self.last_sent_data = None
    
    def init(self):
        try:
            #channel and stub creation
            with grpc.insecure_channel(self.server) as channel:
                stub = edge_pb2_grpc.SensorServiceStub(channel)

                #impacchettamento dei dati
                data = edge_pb2.SensorIdentification(
                    serial_number = self.serial_number
                )

                result = stub.Configuration(data)
                logger.info("Initialization completed")
    
                return result
        except grpc.RpcError as e:
            logger.error("gRPC error: " + str(e))
            traceback.print_exc()
            
        except Exception as e:
            logger.error("Comunication error: " + str(e))
            traceback.print_exc()
            
    def getHourAndMonth(self, timestamp):
        #"2025-05-17T15:33:56.3260074+00:00"
        raw_hour = int(timestamp[11:13])

        if(raw_hour > 4 and raw_hour <= 9):
            hour = 0
        elif(raw_hour > 9 and raw_hour <= 12):
            hour = 1
        elif(raw_hour > 12 and raw_hour <= 17):
            hour = 2
        elif(raw_hour > 17 and raw_hour <= 22):
            hour = 3
        else:
            hour = 4
        
        month = int(timestamp[5:7]) - 1

        return hour, month


    def getData(self):
        timestamp = str(datetime.datetime.now())

        #per ogni metrica, con probabilità dell'1% non viene generato il valore
        temperature = advancedSensorSimulator.get_temperature(timestamp)
        humidity = advancedSensorSimulator.get_humidity(timestamp)
        brightness = advancedSensorSimulator.get_brightness(timestamp)
        air_quality = advancedSensorSimulator.get_air_quality(timestamp)                                                           

        logger.info("Data successfully recordered")

        return {
            "deviceID": self.deviceID,
            "parkID": self.parkID,
            "temperature": temperature,
            "humidity": humidity,
            "brightness": brightness,
            "air_quality": air_quality,
            "timestamp": timestamp
        }
    
    def validateData(self, data):
        logger.info("Validating data")
        for valore in data.values():
            if valore == None:
                return False
        
        if self.last_sent_data is None:
            #i dati non sono stati mai inviati, quindi se sono validi procediamo con l'invio
            return True
        
        #i dati non già stati inviati almeno una volta, quindi dobbiamo verificare se vale la pena inviarli
        for key, threshold in thresholds.items():
            if abs(data[key] - self.last_sent_data[key]) > threshold:
                return True

        return False
    
    def sendData(self, raw_data):
        try:
            #channel and stub creation
            with grpc.insecure_channel(self.server) as channel:
                stub = edge_pb2_grpc.SensorServiceStub(channel)

                #impacchettamento dei dati
                data = edge_pb2.SensorData(
                    deviceID = raw_data["deviceID"],
                    parkID = raw_data["parkID"],
                    temperature = raw_data["temperature"],
                    humidity = raw_data["humidity"],
                    brightness = raw_data["brightness"],
                    air_quality = raw_data["air_quality"],
                    timestamp = raw_data["timestamp"],
                )

                future = stub.SendData.future(data)
                result = future.result()
    
                if result.success:
                    logger.info("Data sent successfully: " + str(result.message))
                else:
                    logger.error("Unsuccessfully comunication: " + str(result.message))
        except grpc.RpcError as e:
            logger.error("gRPC error: " + str(e))
            traceback.print_exc()
            exit(0)
        except Exception as e:
            logger.error("Comunication error: " + str(e))
            traceback.print_exc()
            exit(0)

    
    def execute(self):
        #3) invia i dati
        #4) attesa del turno
        while True:
            try:
                #1) collezione dei dati
                data = self.getData()
                logger.info("Recorded data: " + str(data))

                #2) validazione dei dati
                isValid = self.validateData(data)
                
                #3) invia i dati
                if isValid:
                    logger.info("Valid data, they will be sent")
                    self.sendData(data)
                    self.last_sent_data = data
                else: 
                    logger.info("Invalid data, they will not be sent")

                #4) attesa di {interval} secondi prima di riminciare il ciclo
                logger.info("sleeping " + str(self.interval) + " seconds until next collection")
                time.sleep(int(self.interval))
                                
            except Exception as e:
                logger.error("Problem in the sensor execution: " + str(e))
                traceback.print_exc()
                exit(0)

def main():
    logging.basicConfig(
        filename='edge.log',
        filemode='w',  # append (usa 'w' per sovrascrivere ogni volta)
        level=logging.INFO,
        format='%(asctime)s - %(levelname)s - %(message)s'
    )

    # Apri e leggi il file JSON
    with open("config.json", "r") as f:
        config = json.load(f)
    
    server = config.get("server")
    serial_number = config.get("serial_number")

    logger.info("Starting edge device with serial number " + str(serial_number))

    #istanza del collector
    sensor = ParkSensor(server, serial_number)
    sensor.execute()


if __name__ == "__main__":
    main()