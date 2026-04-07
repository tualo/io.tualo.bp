from io import BytesIO
import tensorflow as tf
from tensorflow import keras
from tensorflow.keras import layers
import base64
import matplotlib.pyplot as plt
from http.server import BaseHTTPRequestHandler, HTTPServer
from PIL import Image

hostName = "localhost"
serverPort = 8080
use_model = '../models/m32x32x150.keras'

class MyServer(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header("Content-type", "text/html")
        self.end_headers()
        self.wfile.write(bytes("<html><head><title>https://pythonbasics.org</title></head>", "utf-8"))
        self.wfile.write(bytes("<p>Request: %s</p>" % self.path, "utf-8"))
        self.wfile.write(bytes("<body>", "utf-8"))
        self.wfile.write(bytes("<p>This is an example web server.</p>", "utf-8"))
        self.wfile.write(bytes("</body></html>", "utf-8"))

    def do_PUT(self):
        # path = self.translate_path(self.path)
        if self.path.endswith('/'):
            self.send_response(405, "Method Not Allowed")
            self.wfile.write("PUT not allowed on a directory\n".encode())
            return
        else:

            length = int(self.headers['Content-Length'])
            content = self.rfile.read(length)
            # print(content)
            im = Image.open(BytesIO(base64.b64decode(content)))
            im = im.convert('L')
            resample = Image.NEAREST
            im = im.resize((image_size, image_size), resample)
            img_array = keras.utils.img_to_array(im )

            # img = keras.utils.load_img(base64.b64decode(content), color_mode='grayscale', target_size=(image_size, image_size), interpolation='nearest')
            # img_array = keras.utils.img_to_array(img)
            img_array = keras.ops.expand_dims(img_array, 0)  # Create batch axis
            predictions = model.predict(img_array)
            prediction_layer = tf.keras.layers.Dense(1, activation='sigmoid')
            score = float(keras.ops.sigmoid(predictions[0][0]))

            self.send_response(200)
            self.send_header("Content-type", "text/html")
            self.end_headers()

            self.wfile.write(bytes(f"{100 * (1 - score):.2f}", "utf-8"))



if __name__ == "__main__":        

    image_size=32
    model = tf.keras.models.load_model(use_model)
    # Check its architecture
    model.summary()

    webServer = HTTPServer((hostName, serverPort), MyServer)
    print("Server started http://%s:%s" % (hostName, serverPort))

    try:
        webServer.serve_forever()
    except KeyboardInterrupt:
        pass

    webServer.server_close()
    print("Server stopped.")




def none():

    # model = tf.keras.models.load_model('saved_model.pb')

    # Check its architecture
    # model.summary()
    # , target_size=(2000, 4000)
    # img = keras.utils.load_img("/Users/thomashoffmann/Documents/Projects/go/create-dataset/10009.4.81.jpg",

    # img = keras.utils.load_img("/Users/thomashoffmann/Documents/Projects/go/create-dataset/10008.2.176.jpg",
    
    img = keras.utils.load_img("/Users/thomashoffmann/Documents/Projects/go/deep-test/data/X/18288.1.1.jpg",
        color_mode='grayscale', 
        target_size=(image_size, image_size),
        interpolation='nearest')

    # plt.imshow(img)

    img_array = keras.utils.img_to_array(img)
    img_array = keras.ops.expand_dims(img_array, 0)  # Create batch axis


    predictions = model.predict(img_array)
    # print(predictions)

    prediction_layer = tf.keras.layers.Dense(1, activation='sigmoid')
    # prediction_batch = prediction_layer(feature_batch_average)
    # print(prediction_batch.shape)

    score = float(keras.ops.sigmoid(predictions[0][0]))
    print(f"This image is {100 * (1 - score):.2f}% X and {100 * score:.2f}% O.")




    img = keras.utils.load_img("/Users/thomashoffmann/Documents/Projects/go/deep-test/data/O/18300.0.0.jpg",
        color_mode='grayscale', 
        target_size=(image_size, image_size),
        interpolation='nearest')

    # plt.imshow(img)

    img_array = keras.utils.img_to_array(img)
    img_array = keras.ops.expand_dims(img_array, 0)  # Create batch axis


    predictions = model.predict(img_array)
    # print(predictions)

    prediction_layer = tf.keras.layers.Dense(1, activation='sigmoid')
    # prediction_batch = prediction_layer(feature_batch_average)
    # print(prediction_batch.shape)

    score = float(keras.ops.sigmoid(predictions[0][0]))
    print(f"This image is {100 * (1 - score):.2f}% X and {100 * score:.2f}% O.")

    # print(f'input_layer_name={model.input.name}')
    # output_layer_name = model.output.name.split(':')[0]
    # print(f'output_layer_name={output_layer_name}')
    # plt.show()