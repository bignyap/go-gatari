from diagrams import Cluster, Diagram, Edge
from diagrams.aws.network import APIGateway
from diagrams.onprem.client import User
from diagrams.onprem.compute import Server
from diagrams.onprem.database import PostgreSQL
from diagrams.onprem.inmemory import Redis
from diagrams.generic.compute import Rack
from diagrams.programming.language import Go
from diagrams.generic.device import Mobile

graph_attr = {
    "fontsize": "22",  # Fallback title size
    "labelloc": "t",  # Title at top
    "labeljust": "c",  # Centered title
    "fontname": "Helvetica-Bold",
    "labelfontname": "Helvetica-Bold",
    "labelfontsize": "28",
    "bgcolor": "white",
    "pad": "0.5",
    "ranksep": "1.0",
}

node_attr = {
    "fontsize": "14",
    "fontname": "Helvetica",
    "style": "filled",
    "fillcolor": "#f5f5f5",
    "color": "#2c3e50",
    "fontcolor": "#2c3e50",
}

edge_attr = {
    "fontsize": "12",
    "fontname": "Helvetica",
    "fontcolor": "#34495e",
    "color": "#7f8c8d",
}

with Diagram("GATARI Platform Architecture", show=False, filename="gatari_full_architecture", direction="LR",
             graph_attr=graph_attr, node_attr=node_attr, edge_attr=edge_attr):
    
    end_user = User("End User")
    admin_user = User("Administrator")

    with Cluster("Gateway Layer"):
        api_gateway = APIGateway("API Gateway")

    with Cluster("Application Layer"):
        application = Server("Application Service")

    with Cluster("GATARI Platform"):

        with Cluster("Gatekeeper Service"):
            gatekeeper = Go("Gatekeeper")
            local_cache = Rack("Local Cache\n(1s TTL)")
            redis_cache = Redis("Redis Cache\n(2s TTL)")

        with Cluster("Admin Service"):
            admin_ui = Mobile("Admin Interface")
            admin_backend = Go("Admin Backend")

        with Cluster("Database Layer"):
            database = PostgreSQL("PostgreSQL")

        with Cluster("Messaging Layer"):
            redis_pubsub = Redis("Redis PubSub")

    # Admin Flow
    admin_user >> admin_ui >> admin_backend
    admin_backend >> Edge(label="Reads/Writes", color="#2980b9", style="bold") >> database
    admin_backend >> Edge(label="Publishes Updates", color="#2980b9", style="dashed") >> redis_pubsub
    redis_pubsub >> Edge(label="Config Events", color="#2980b9", style="dotted") >> gatekeeper

    # Gatekeeper Cache Chain
    gatekeeper >> Edge(label="Read-Through", color="#27ae60") >> local_cache
    local_cache >> Edge(label="Fallback", color="#27ae60", style="dashed") >> redis_cache
    redis_cache >> Edge(label="DB Miss", color="#27ae60", style="dotted") >> database

    # User Request Flow
    end_user >> Edge(label="API Request", color="#8e44ad", style="bold") >> api_gateway
    api_gateway >> gatekeeper >> application
    application >> Edge(label="Response", color="#8e44ad") >> gatekeeper >> api_gateway