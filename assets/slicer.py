from PIL import Image
import os
import sys


def make_black_pixels_transparent(image):
    # Convert the image to RGBA mode
    image = image.convert("RGBA")

    # Get the pixel data
    pixel_data = image.getdata()

    # Create a new pixel list with black pixels made transparent
    new_pixel_data = []
    for pixel in pixel_data:
        if pixel[:3] == (0, 0, 0):  # Check if the pixel is completely black
            new_pixel_data.append((0, 0, 0, 0))  # Make it transparent
        else:
            new_pixel_data.append(pixel)

    # Update the image with the new pixel data
    image.putdata(new_pixel_data)

    return image


def split_sprite_sheet(input_image_path, output_folder):
    # Open the sprite sheet image
    sprite_sheet = Image.open(input_image_path)

    # Get the dimensions of each sprite (assuming they are all 32x32)
    sprite_width, sprite_height = 32, 32

    # Get the number of rows and columns in the sprite sheet
    num_rows = sprite_sheet.height // sprite_height
    num_cols = sprite_sheet.width // sprite_width

    # Create the output folder if it doesn't exist
    os.makedirs(output_folder, exist_ok=True)

    # Loop through the sprite sheet and extract each sprite
    for row in range(num_rows):
        for col in range(num_cols):
            # Define the coordinates of the current sprite
            left = col * sprite_width
            top = row * sprite_height
            right = left + sprite_width
            bottom = top + sprite_height

            # Crop and make black pixels transparent
            sprite = sprite_sheet.crop((left, top, right, bottom))
            sprite = make_black_pixels_transparent(sprite)

            # Calculate the zero-indexed row number and column number (frame)
            sprite_index = row
            sprite_frame = col

            # Save the sprite as a new image
            sprite.save(
                f"{output_folder}/sprite_{sprite_index}_{sprite_frame}.png")


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("specify input image")
        sys.exit(1)

    input_image_path = sys.argv[1]
    output_folder = "sprites"

    split_sprite_sheet(input_image_path, output_folder)
