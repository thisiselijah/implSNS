import { PostCard } from "@/components/card";

export default function Feed({ FeedData = [] }) {
    console.log("Feed component rendered with posts:", FeedData);
    return (
        <div className="">
            {FeedData.map((post) => (
                <PostCard key={post.post_id || post.id} post={post} />
            ))}
            
        </div>
    );
}