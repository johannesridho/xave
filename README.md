# Xave
Xave makes cities save by detecting suspicious activities (vehicle and fire accidents, road conflicts like people fighting) using CCTVs and send FB Message to city leaders or police

## Links
https://www.facebook.com/Xave-2595813077099497

## Inspiration
More than half of the world’s population now live in urban areas. By 2050, that figure will have risen to 6.5 billion people – two-thirds of all humanity. Sustainable development cannot be achieved without significantly transforming the way we build and manage our urban spaces. Problems need to be addressed as fast as possible. Here I tried to build a solution to report suspicious activities around cities, for example, vehicle or fire accidents and violence activities in cities.

## What it does
Detect suspicious activities from CCTVs videos (for example, vehicle or fire accidents and violence activities) and report to relevant users (city leaders or police) with Facebook Messenger

## How I built it
I created Facebook page and Messenger. Then I created several Amazon Lambda functions to check if there is a new video uploaded to Amazon S3, then do several preprocessing tasks and use Amazon Rekognition to do the computer vision part. Then do several other tasks and send the result to Facebook Messenger.

## User's privacy
Every report generated will only be sent to relevant city leaders. City leaders can only access the relevant data from their CCTV. 

## Potential impact
When this product installed in CCTVs, city leaders and police can address city issues faster, without the need to manually check CCTVs or wait to get complaints from the citizens, or for example, wait for a building to get burned into ashes.

## Diversity and inclusion
Crime rate faced people with disability is 2.5 times more than the others. Up to 90% of women with disability have been sexually assaulted. Xave helps the assaulted people with disability by notifying the police faster, which means polices can find and help them faster. This is very significant especially for visually impaired people, which are often don't know who assaulted them. 

## What's next for Xave
Add more features, there are so much information we can get from cities CCTVs 
