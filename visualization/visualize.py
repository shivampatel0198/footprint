import json

import matplotlib
matplotlib.use('Agg')

from matplotlib import pyplot as plt
import matplotlib.animation as animation

import numpy as np

# (0) Set up figure
fig = plt.figure()

# (1) Read JSON data into a map
with open('data.json') as file:
    data = json.load(file)

# (2) Plot points on a grid
frames = []
for timestep in range(len(data)):

    # Pull coordinates from all nodes in this timestep
    coords = np.array([
        [data[timestep][node]['Loc']['X'], 
        data[timestep][node]['Loc']['Y'], 
        'red' if (data[timestep][node]['Infected']) else 'black']
        for node in range(len(data[timestep]))
    ])

    plt.axis('off')
    x, y, c = coords.T
    frames.append([plt.scatter(x=x,y=y,c=c)])

ani = animation.ArtistAnimation(fig, frames, interval=1000, blit=False,
                                repeat=False)

writer = animation.FFMpegWriter(bitrate=6000000, fps=10)
ani.save('viz.mp4', writer=writer)