  /*
 * HTTP Client GET Request
 * Copyright (c) 2018, circuits4you.com
 * All rights reserved.
 * https://circuits4you.com 
 * Connects to WiFi HotSpot. */

#include <WiFi.h>;
#include <HTTPClient.h>;

/* Set these to your desired credentials. */
const char *ssid = "Nakula_Ext";  //ENTER YOUR WIFI SETTINGS
const char *password = "#Bjm01bsd#";

//Web/Server address to read/write from 
const char *host1 = "http://192.168.18.157:8080/iot";   //https://circuits4you.com website or IP address of server
const char *host2 = "http://192.168.18.157:8080/iot2"; 

//pin ultrasonic sensor
const int trigPin[] = {22, 4, 26, 18};
const int echoPin[] = {23, 2, 27, 19};

//ultrasonic variable
float SOUND_VELOCITY = 0.034;
long duration;
float distanceCm[4];
float distanceCmStandard[4];
String type[] = {"0", "1", "2", "3", "4"};
String cond = "0";

//=======================================================================
//                    Power on setup
//=======================================================================

void setup() {
  pinMode(trigPin[0], OUTPUT);
  pinMode(echoPin[0], INPUT);
  for (int i = 1; i < 4; i++){
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

void send_post(String httpRequestData){
  //harapannya ketika sudah ke detect barang ada yang masuk, 
  WiFiClient wifiClient;
  HTTPClient http;    //Declare object of class HTTPClient
  String Link = host2;
  if(httpRequestData == "0" || httpRequestData == "1"){
    Link = host1;
  }
  
  http.begin(wifiClient, Link);     //Specify request destination
  
  http.addHeader("Content-Type", "text/plain");
  int httpCode = http.POST(httpRequestData);
  String payload = http.getString();    //Get the response payload

  Serial.println(httpCode);   //Print HTTP return code
  Serial.println(payload);    //Print request response payload

  http.end();  //Close connection
}

//=======================================================================
//                    Main Program Loop
//=======================================================================
void loop() {
  // Clears the trigPin
  digitalWrite(trigPin[0], LOW);
  delayMicroseconds(2);
  // Sets the trigPin on HIGH state for 10 micro seconds
  digitalWrite(trigPin[0], HIGH);
  delayMicroseconds(10);
  digitalWrite(trigPin[0], LOW);

  // Reads the echoPin, returns the sound wave travel time in microseconds
  duration = pulseIn(echoPin[0], HIGH);
  
  // Calculate the distance
  distanceCm[0] = duration * SOUND_VELOCITY /2;
  if(distanceCm[0] < 100){
    if(cond == "0"){
      cond = "1";
      Serial.print("human ");
      Serial.println(cond);
      send_post(type[1]);
    }
  }else{
    if(cond == "1"){
      cond = "0";
      Serial.print("human ");
      Serial.println(cond);
      send_post(type[0]);
    }
  }

  for (int i=1; i < 4; i++){
    // Clears the trigPin
    Serial.println(i);
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
    Serial.println(distanceCm[i]);

    if(abs(distanceCm[i] - distanceCmStandard[i]) > 1){
      Serial.println("sampah");
      send_post(type[i+1]);
    }
  }
  delay(1000);
}