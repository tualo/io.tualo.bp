import tensorflow as tf
from tensorflow import keras
from tensorflow.keras import layers
from tensorflow.keras.preprocessing.image import ImageDataGenerator
import matplotlib.pyplot as plt

data_dir = './data'
img_height=32
img_width=32
batch_size=50
epochs= 150

train_ds = tf.keras.utils.image_dataset_from_directory(
    data_dir,
    validation_split=0.1,
    subset="training",
    color_mode='grayscale',
    seed=123,
    image_size=(img_height, img_width),
    batch_size=batch_size)

val_ds = tf.keras.utils.image_dataset_from_directory(
    data_dir,
    validation_split=0.1,
    color_mode='grayscale',
    subset="validation",
    seed=123,
    image_size=(img_height, img_width),
    batch_size=batch_size)

modelX = keras.Sequential([
    keras.Input(shape=(img_height, img_width, 1)),
    layers.Conv2D(16, 3, padding='same'),
    layers.Conv2D(32, 3, padding='same'),
    layers.MaxPooling2D(pool_size=(2, 2)),
    layers.Flatten(),
    layers.Dense(10)
])

num_classes = 5

model = tf.keras.Sequential([
  tf.keras.layers.Rescaling(1./255),
  tf.keras.layers.Conv2D(32, 3, activation='relu'),
  tf.keras.layers.MaxPooling2D(),
  tf.keras.layers.Conv2D(32, 3, activation='relu'),
  tf.keras.layers.MaxPooling2D(),
  tf.keras.layers.Conv2D(32, 3, activation='relu'),
  tf.keras.layers.MaxPooling2D(),
  tf.keras.layers.Flatten(),
  tf.keras.layers.Dense(128, activation='relu'),
  tf.keras.layers.Dense(num_classes)
])

model.compile(optimizer='adam',
                loss=tf.losses.SparseCategoricalCrossentropy(from_logits=True),
                metrics=['accuracy'])

model.fit(train_ds, validation_data=val_ds, epochs=epochs)
buf = "models/m%dx%dx%d.keras" % (img_width, img_height, epochs)
model.save(buf)
# tf.saved_model.save(model, "./")
model.summary()

# print(f'input_layer_name={model.input.name}')
# output_layer_name = model.output.name.split(':')[0]
# print(f'output_layer_name={output_layer_name}')

# image_batch, label_batch = next(iter(train_ds))

# plt.figure(figsize=(10, 10))
# for i in range(9):
#  ax = plt.subplot(3, 3, i + 1)
#  plt.imshow(image_batch[i].numpy().astype("uint8"))
#  label = label_batch[i]
  # plt.title(class_names[label])
#  plt.axis("off")

