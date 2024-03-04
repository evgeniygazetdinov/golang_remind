#include <OneWire.h>
#include <DallasTemperature.h>
#include <ESP8266WiFi.h>
#include <ESP8266WebServer.h>
#include <WiFiClient.h>

#define ONE_WIRE_BUS 2  // NodeMCU GPIO pin number where DS18b20 data pin connected

OneWire oneWire(ONE_WIRE_BUS);
DallasTemperature sensors(&oneWire);
ESP8266WebServer server(80);

//wifi settings
const char* ssid = "TP-Link_AB03"; //Use your own SSID
const char* password =  "";  //Use your own Passoword

void setup(void)
{
  // start serial port
  Serial.begin(9600);
  sensors.begin();

  //Wi-Fi 
  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    delay(2000);
    Serial.println("Connecting ...");
  }

  Serial.print("Connected to WiFi, IP Address is: ");
  Serial.println(WiFi.localIP()); 
  server.on("/", Connect);
  server.onNotFound(Err_Connect);
  server.begin();
  
}


int room_temp;

void loop(void)
{ 
  sensors.requestTemperatures();
  room_temp = sensors.getTempCByIndex(0);
  room_temp += 20;
  if(room_temp != DEVICE_DISCONNECTED_C) 
  {

  } 
  
  else
  {
    Serial.println("Error: Could not read temperature !");
  }
   server.handleClient();
}

void Connect() {
  
  server.send(200, "text/json", HTML_Code(room_temp));
}

void Err_Connect(){
  server.send(200, "text/json", "{\"temp\": \"\"\"}");
}

String HTML_Code(int Temp){
  String message = "";
  message += "{\"temp\":\"";
  message += Temp;
  message +="\"}";


  return message;
}