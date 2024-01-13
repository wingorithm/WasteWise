/*
 * HTTP Client GET Request
 * Copyright (c) 2018, circuits4you.com
 * All rights reserved.
 * https://circuits4you.com 
 * Connects to WiFi HotSpot. */

#include <ESP8266WiFi.h>
#include <WiFiClient.h> 
#include <ESP8266WebServer.h>
#include <ESP8266HTTPClient.h>

/* Set these to your desired credentials. */
const char *ssid = "ADRIEL";  //ENTER YOUR WIFI SETTINGS
const char *password = "cintakaadriel";

//Web/Server address to read/write from 
const char *host = "192.168.18.13:9000";   //https://circuits4you.com website or IP address of server

//pin ultrasonic sensor
const int trigPin[] = {D1, D3, D5, D7};
const int echoPin[] = {D2, D4, D6, D8};

//ultrasonic variable
#define SOUND_VELOCITY 0.034
long duration;
float distanceCm[4];
float distanceCmStandard[4];
String type[] = {"Recycle", "Organic", "Another"}

//=======================================================================
//                    Power on setup
//=======================================================================

void setup() {
  for (int i = 0; i < 4; i++){
    pinMode(trigPin[i], OUTPUT);
    pinMode(echoPin[i], INPUT);

    // define first distance
    // Clears the trigPin
    digitalWrite(trigPin[i], LOW);
    delayMicroseconds(2);
    // Sets the trigPin on HIGH state for 10 micro seconds
    digitalWrite(trigPin[i], HIGH);
    delayMicroseconds(10);
    digitalWrite(trigPin[i], LOW);

    // Reads the echoPin, returns the sound wave travel time in microseconds
    duration = pulseIn(echoPin[i], HIGH);
    
    // Calculate the distance
    distanceCmStandard[i] = duration * SOUND_VELOCITY/2;
  }
  delay(1000);
  Serial.begin(9600);
  WiFi.mode(WIFI_OFF);        //Prevents reconnection issue (taking too long to connect)
  delay(1000);
  WiFi.mode(WIFI_STA);        //This line hides the viewing of ESP as wifi hotspot
  
  WiFi.begin(ssid, password);     //Connect to your WiFi router
  Serial.println("");

  Serial.print("Connecting");
  // Wait for connection
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  //If connection successful show IP address in serial monitor
  Serial.println("");
  Serial.print("Connected to ");
  Serial.println(ssid);
  Serial.print("IP address: ");
  Serial.println(WiFi.localIP());  //IP address assigned to your ESP
}

//=======================================================================
//                    Main Program Loop
//=======================================================================
void loop() {
  for (int i=0; i < 4; i++){
    // Clears the trigPin
    digitalWrite(trigPin[i], LOW);
    delayMicroseconds(2);
    // Sets the trigPin on HIGH state for 10 micro seconds
    digitalWrite(trigPin[i], HIGH);
    delayMicroseconds(10);
    digitalWrite(trigPin[i], LOW);

    // Reads the echoPin, returns the sound wave travel time in microseconds
    duration = pulseIn(echoPin[i], HIGH);
    
    // Calculate the distance
    distanceCm[i] = duration * SOUND_VELOCITY/2;

    if(abs(distanceCm[i] - distanceCmStandard[i]) > 1){
      WiFiClient wifiClient;
      HTTPClient http;    //Declare object of class HTTPClient
      String Link = "http://192.168.18.13:9000";
      
      http.begin(wifiClient, Link);     //Specify request destination
      
      http.addHeader("Content-Type", "application/json");
      String httpRequestData = "{\"type\":\""+type[i]+"\"}";
      int httpCode = http.POST(httpRequestData);
      String payload = http.getString();    //Get the response payload

      Serial.println(httpCode);   //Print HTTP return code
      Serial.println(payload);    //Print request response payload

      http.end();  //Close connection
    }
  }
  
  delay(10);
}