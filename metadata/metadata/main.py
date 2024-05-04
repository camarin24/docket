def main():
    print("Hello, World!")
    print("This is the main function.")
    # Get files from db and check if they metadata must be extracted
    # Create a chain of extractors and processors to extract metadata from files
    # Create a unique folder structure for each file to store extracted metadata such as thumbnails, raw contents, and metadata
    #    - Put all metadata in a single json file may be a good idea
    # Send metadata to a vector db for indexing and searching
    # Save metadata in a database for indexing and searching
    #    - Figure out how to store some chunks of metadata in the database for indexing and searching
    # Update the file in the database
    # This api must be able to being called from some other api to extract metadata for a specific file

