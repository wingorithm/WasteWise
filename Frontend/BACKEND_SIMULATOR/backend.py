from flask import Flask, jsonify
from flask_cors import CORS

import time
import threading

app = Flask(__name__)
CORS(app)


@app.route('/idle', methods=['GET'])
def long_polling_idle():
    timeout = 12  # Response delay in seconds

    # Simulate a long process
    time.sleep(timeout)

    return jsonify({"detected": True})

@app.route('/classify', methods=['GET'])
def classify():
    timeout = 60  # Response delay in seconds

    # Simulate a long process
    time.sleep(timeout)

    return jsonify({"type": 0})

@app.route('/voice', methods=['GET'])
def long_polling():
    timeout = 15  # Response delay in seconds

    # Simulate a long process
    time.sleep(timeout)

    return jsonify({"name": "Erwin G"})

if __name__ == '__main__':
    # Run the Flask app in a separate thread
    threading.Thread(target=app.run, kwargs={'debug': True, 'use_reloader': False}).start()