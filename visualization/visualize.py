import json

import matplotlib
matplotlib.use('Agg')

from matplotlib import pyplot as plt
import matplotlib.animation as animation

import numpy as np

# VISUALIZATION CONSTANTS
infected_color     = 'red'
not_infected_color = 'black'
stepsize           = 2

# (0) Set up figure
fig, ax = plt.subplots()
fig.subplots_adjust(left=0, bottom=0, right=1, top=1, wspace=None, hspace=None)
plt.autoscale(enable=True)

ax.grid(color='gray', alpha=0.3, linestyle='dashed', linewidth=0.2, zorder=-1)

plt.xticks(np.arange(-100, 100, stepsize))
plt.yticks(np.arange(-100, 100, stepsize))


# (1) Read JSON data into a map
with open('data.json') as file:
    data = json.load(file)

# (2) Plot points on a grid
frames = []
for timestep in range(len(data)):

    # Pull coordinates from all nodes and assign colors
    coords = np.array([
        [data[timestep][node]['Loc']['X'], 
        data[timestep][node]['Loc']['Y'], 
        infected_color if (data[timestep][node]['Infected']) else not_infected_color]
        for node in range(len(data[timestep]))
    ])

    x, y, c  = coords.T
    frames.append([plt.scatter(x=x,y=y,c=c, alpha=0.5)])

ani = animation.ArtistAnimation(fig, frames, interval=1000, blit=False,
                                repeat=False)

# Play with these values (bitrate, fps, and dpi) for higher quality video
writer = animation.FFMpegWriter(bitrate=2000, fps=8)
ani.save('viz.mp4', writer=writer, dpi=500)