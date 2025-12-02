import requests
from bs4 import BeautifulSoup, Tag
import re
import json
import time
from typing import List, Dict, Optional, Set, Any, Tuple

# --- Constants ---
BASE_URL = "https://whatsonzwift.com"
LISTING_URL = BASE_URL + "/search?sport=bike&d=all&sp=all&l=all&z=all&k=&s=new&o[zc]=1&o[zw]=1&o[zf]=1&o[c]=1&page={page}#results"
REQUEST_DELAY_SECONDS = 0.5
MAX_PAGES = 129 

# --- Data Structures (for clarity) ---
WorkoutStep = Dict[str, Any]
WorkoutData = Dict[str, Any]

def get_soup(url: str) -> Optional[BeautifulSoup]:
    """Fetches a URL and returns a BeautifulSoup object, handling errors."""
    try:
        response = requests.get(url, timeout=10)
        response.raise_for_status()
        return BeautifulSoup(response.text, 'html.parser')
    except requests.exceptions.RequestException as e:
        return None

def parse_time_duration(time_str: str) -> int:
    """Parse time string like '10min', '30sec', or '4hr 10min' to seconds."""
    total_seconds = 0
    # Match for hours (e.g., '4hr')
    h_match = re.search(r'(\d+)\s*hr', time_str)
    if h_match:
        total_seconds += int(h_match.group(1)) * 3600
        
    # Match for minutes (e.g., '10min')
    m_match = re.search(r'(\d+)\s*min', time_str)
    if m_match:
        total_seconds += int(m_match.group(1)) * 60
        
    # Match for seconds (e.g., '30sec')
    s_match = re.search(r'(\d+)\s*sec', time_str)
    if s_match:
        total_seconds += int(s_match.group(1))
        
    return total_seconds

def parse_step(line: str) -> Optional[WorkoutStep]:
    """
    Converts individual step segments into a dict, guaranteeing start_power and end_power.
    Uses robust regex to handle inconsistent spacing.
    """
    # 1. Robust Pattern for RAMP: {time} from{optional_space}{power1}to{optional_space}{power2}% FTP
    ramp = re.match(r"(.+)\s*from\s*(\d{1,3})\s*to\s*(\d{1,3})% FTP", line)
    
    # 2. Robust Pattern for STEADY: {time} {optional_space}@{optional_space}{power}% FTP
    steady = re.match(r"(.+)\s*@\s*(\d{1,3})% FTP", line)
    
    if ramp:
        duration_sec = parse_time_duration(ramp.group(1))
        start_power = int(ramp.group(2))
        end_power = int(ramp.group(3))
        if duration_sec > 0:
            return {
                "duration": duration_sec, 
                "start_power": start_power, 
                "end_power": end_power, 
                "type": "ramp"
            }
            
    if steady:
        duration_sec = parse_time_duration(steady.group(1))
        power = int(steady.group(2))
        if duration_sec > 0:
            # FIX: For steady state, set start and end power to the same value for consistency
            return {
                "duration": duration_sec, 
                "start_power": power,
                "end_power": power, 
                "type": "steady"
            }
            
    # Handle the 'Free Ride' block (warmup/cooldown/free ride without explicit power listed)
    duration_sec = parse_time_duration(line)
    if duration_sec > 0 and "% FTP" not in line:
         # Assign a default Z2 power for general riding steps
        default_power = 60 
        return {
            "duration": duration_sec, 
            "start_power": default_power, 
            "end_power": default_power, 
            "type": "steady" 
        }
        
    return None

def parse_repetition(line: str) -> Tuple[int, str]:
    """Extracts repetition count (e.g., '3x') and the cleaned remaining step string."""
    rep_match = re.match(r'(\d+)\s*x\s*(.*)', line, re.I)
    if rep_match:
        count = int(rep_match.group(1))
        content = rep_match.group(2).strip()
        return count, content
    return 1, line

def extract_steps(article: Tag) -> List[WorkoutStep]:
    """
    Extracts workout steps by targeting the 'textbar' div elements, 
    handling repetitions and multi-segment lines.
    """
    all_steps: List[WorkoutStep] = []
    
    step_elements = article.find_all('div', class_='textbar')
    
    for el in step_elements:
        line = el.get_text(strip=True).replace('\xa0', ' ').strip()
        
        # 1. Check for repetitions and separate the count from the content
        reps, content = parse_repetition(line)
        
        # 2. Split content into individual segments (intervals separated by comma)
        # Handle the case where content is empty if the line was just "3x" (shouldn't happen here)
        if not content:
            continue
            
        segments = [s.strip() for s in content.split(',')]
        
        parsed_segments: List[WorkoutStep] = []

        # 3. Parse each segment into a WorkoutStep dictionary
        for segment in segments:
            step = parse_step(segment)
            if step:
                parsed_segments.append(step)

        # 4. If segments were successfully parsed, apply the repetition factor
        if parsed_segments:
            for _ in range(reps):
                all_steps.extend(parsed_segments)
        # Optional: Keep the warning if a line that looks like a step couldn't be parsed
        elif ("min" in line or "sec" in line or "hr" in line) and "% FTP" in line:
             print(f'    [warning] Unparsed line (check spacing/format): {line}')

    return all_steps

# --- Remaining functions (process_workout, find_workout_links, etc.) unchanged from previous step ---

def process_workout(workout_url: str) -> Optional[WorkoutData]:
    """Fetches a single workout page, extracts its name and steps."""
    print(f"  ... Processing {workout_url}")
    soup = get_soup(workout_url)
    if not soup:
        return None

    article = soup.find('article')
    if not article or not isinstance(article, Tag):
        print(f'    [error] Could not find article content on {workout_url}')
        return None

    # Find the first heading element
    name_elements = article.find_all(['h1', 'h2', 'h3'], limit=1)
    
    if name_elements:
        name_el = name_elements[0]
        name = name_el.get_text(strip=True)
    else:
        name = workout_url.split('/')[-1]

    steps = extract_steps(article)
    
    # Check if a known workout has steps (useful for debugging)
    if not steps and name:
        pass 

    return {
        "name": name,
        "url": workout_url,
        "steps": steps
    }

def find_workout_links(soup: BeautifulSoup) -> List[str]:
    """Returns all unique full workout detail links on one search page."""
    links: List[str] = []
    
    # Target <a> tags that have an href starting with the full base URL followed by '/workouts/'
    workout_anchors = soup.find_all(
        'a', 
        href=re.compile(r'^https://whatsonzwift.com/workouts/') 
    )

    for a in workout_anchors:
        href = a.get('href')
        if href:
            links.append(str(href)) 
            
    return list(set(links))

def scrape_all_workouts(max_pages: int = MAX_PAGES) -> List[WorkoutData]:
    """Main function to iterate through search pages and scrape individual workouts."""
    all_workouts: List[WorkoutData] = []
    processed_urls: Set[str] = set()

    for page in range(1, max_pages + 1):
        print(f"\n--- Scraping Page {page}/{max_pages} ---")
        listing_url = LISTING_URL.format(page=page)
        
        soup = get_soup(listing_url)
        if not soup:
            continue

        workout_links = find_workout_links(soup)
        
        if not workout_links and page > 1:
             print("Reached end of result pages early.")
             break
        
        print(f"Found {len(workout_links)} workout links on this page.")

        for wurl in workout_links:
            if wurl in processed_urls:
                continue
            
            workout_data = process_workout(wurl)
            
            if workout_data:
                all_workouts.append(workout_data)
                processed_urls.add(wurl)
            
            time.sleep(REQUEST_DELAY_SECONDS)

    return all_workouts

def save_results(data: List[WorkoutData], filename: str = "zwift_workouts.json"):
    """Saves the scraped workout data to a JSON file."""
    try:
        with open(filename, "w", encoding="utf-8") as f:
            json.dump(data, f, ensure_ascii=False, indent=2)
        print(f"\n✅ Successfully saved {len(data)} unique workouts to {filename}.")
    except IOError as e:
        print(f"Error saving file: {e}")

# --- Main Execution ---
if __name__ == "__main__":
    workouts = scrape_all_workouts(max_pages=MAX_PAGES)
    save_results(workouts)
