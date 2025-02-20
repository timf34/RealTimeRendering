from PIL import Image
import numpy as np

def create_checkerboard(size=1000, square_size=100):
    """
    Create a black and white checkerboard pattern.
    
    Args:
        size (int): The width and height of the final image in pixels
        square_size (int): The size of each square in the checkerboard
    
    Returns:
        PIL.Image: The generated checkerboard image
    """
    # Calculate number of squares in each row/column
    num_squares = size // square_size
    
    # Create a numpy array for the pattern
    pattern = np.zeros((size, size), dtype=np.uint8)
    
    # Fill the pattern with alternating black and white squares
    for i in range(num_squares):
        for j in range(num_squares):
            if (i + j) % 2 == 0:
                # Fill white square (255 is white in grayscale)
                pattern[i*square_size:(i+1)*square_size, 
                       j*square_size:(j+1)*square_size] = 255
    
    # Convert numpy array to PIL Image
    image = Image.fromarray(pattern, mode='L')
    
    return image

def main():
    # Create the checkerboard
    checkerboard = create_checkerboard(size=1000, square_size=100)
    
    # Save the image
    output_filename = 'checkerboard.png'
    checkerboard.save(output_filename)
    print(f"Checkerboard image saved as {output_filename}")

if __name__ == "__main__":
    main()